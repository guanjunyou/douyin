package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/mq"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"sync"
)

type RelationServiceImpl struct {
	Logger *logrus.Logger
}

// FollowUser 关注用户
func (relationServiceImpl RelationServiceImpl) FollowUser(userId int64, toUserId int64, actionType int) error {
	relationServiceImpl.Logger.Info("FollowUser\n")
	if userId == toUserId {
		return fmt.Errorf("你不能关注(或者取消关注)自己")
	}
	followData := models.FollowMQToUser{
		UserId:       userId,
		FollowUserId: toUserId,
		ActionType:   actionType,
	}
	message, err := json.Marshal(followData)
	if err != nil {
		return err
	}
	mq.FollowRMQ.Publish(message)
	return nil
}

// GetFollows 查询关注列表
func (relationServiceImpl RelationServiceImpl) GetFollows(userId int64) ([]models.User, error) {
	relationServiceImpl.Logger.Info("GetFollows\n")
	var users []models.User
	err := utils.GetMysqlDB().Table("follow").Where("user_id = ? AND is_deleted != ?", userId, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}

	//协程并发更新，isFollow 为 True 前端才能显示已关注
	var wg sync.WaitGroup
	for i := 0; i < len(users); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			users[i].IsFollow = true
		}(i)
	}

	wg.Wait()

	return users, nil
}

// GetFollowers 查询粉丝列表
func (relationServiceImpl RelationServiceImpl) GetFollowers(userId int64) ([]models.User, error) {
	relationServiceImpl.Logger.Info("GetFollowers")
	var users []models.User
	err := utils.GetMysqlDB().Table("follow").Where("follow_user_id = ? AND is_deleted != ?", userId, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetFriends 查询好友列表
func (relationServiceImpl RelationServiceImpl) GetFriends(userId int64) ([]models.User, error) {
	follows, err := relationServiceImpl.GetFollows(userId)
	if err != nil {
		return nil, err
	}
	followers, err := relationServiceImpl.GetFollowers(userId)
	if err != nil {
		return nil, err
	}
	var friends []models.User
	for _, user := range followers {
		if containsID(follows, user.Id) {
			friends = append(friends, user)
		}
	}
	return friends, nil
}

// containsID 辅助函数，用于检查指定的 id 是否在数组中存在
func containsID(arr []models.User, id int64) bool {
	for _, u := range arr {
		if u.Id == id {
			return true
		}
	}
	return false
}

func MakeFollowGroutine(count int) {
	for i := 0; i < count; i++ {
		go mq.FollowRMQ.Consumer()
	}
}

// GetUserFollowing 获取某个用户的关注列表
func GetUserFollowing(userID int64) ([]int64, error) {
	client := utils.GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%d", config.FollowSetKey, userID)

	// 尝试从 Redis 获取数据
	following, err := client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 如果 Redis 中没有数据，则从 MySQL 获取
	if err == redis.Nil {
		following, err = getFollowingFromMySQL(userID)
		if err != nil {
			return nil, err
		}

		// 启动线程异步将数据添加到 Redis
		go func() {
			err := addFollowingToRedis(userID, following)
			if err != nil {
				// 处理错误，例如记录日志
			}
		}()
	}

	// 转换关注列表中的用户ID为int64类型
	var followingIDs []int64
	for _, f := range following {
		var followeeID int64
		_, err := fmt.Sscanf(f, "%d", &followeeID)
		if err != nil {
			return nil, err
		}
		followingIDs = append(followingIDs, followeeID)
	}

	return followingIDs, nil
}

// GetUserFollowers 获取某个用户的粉丝列表
func GetUserFollowers(userID int64) ([]int64, error) {
	client := utils.GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%d", config.FollowerSetKey, userID)

	// 尝试从 Redis 获取数据
	followers, err := client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 如果 Redis 中没有数据，则从 MySQL 获取
	if err == redis.Nil {
		followers, err = getFollowersFromMySQL(userID)
		if err != nil {
			return nil, err
		}

		// 启动线程异步将数据添加到 Redis
		go func() {
			err := addFollowersToRedis(userID, followers)
			if err != nil {
				// 处理错误，例如记录日志
			}
		}()
	}

	// 转换粉丝列表中的用户ID为int64类型
	var followerIDs []int64
	for _, f := range followers {
		var followerID int64
		_, err := fmt.Sscanf(f, "%d", &followerID)
		if err != nil {
			return nil, err
		}
		followerIDs = append(followerIDs, followerID)
	}

	return followerIDs, nil
}

// GetUserFriends 获取某个用户的朋友列表（关注和粉丝的交集）
func GetUserFriends(userID int64) ([]int64, error) {
	client := utils.GetRedisDB()
	ctx := context.Background()

	// 获取关注列表和粉丝列表的键名
	followingKey := fmt.Sprintf("%v%d", config.FollowSetKey, userID)
	followersKey := fmt.Sprintf("%v%d", config.FollowerSetKey, userID)

	// 尝试从 Redis 获取数据
	friends, err := client.SInter(ctx, followingKey, followersKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 如果 Redis 中没有数据，则从 MySQL 获取
	if err == redis.Nil {
		friends, err = getFriendsFromMySQL(userID)
		if err != nil {
			return nil, err
		}

		// 启动线程异步将数据添加到 Redis
		go func() {
			err := addFriendsToRedis(userID, friends)
			if err != nil {
				// 处理错误，例如记录日志
			}
		}()
	}

	// 转换朋友列表中的用户ID为int64类型
	var friendIDs []int64
	for _, f := range friends {
		var friendID int64
		_, err := fmt.Sscanf(f, "%d", &friendID)
		if err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, friendID)
	}

	return friendIDs, nil
}

// 从 MySQL 获取关注列表
func getFollowingFromMySQL(userID int64) ([]string, error) {
	db := utils.GetMysqlDB()

	var follows []models.Follow
	err := db.Where("follower_id = ?", userID).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	var following []string
	for _, follow := range follows {
		following = append(following, fmt.Sprintf("%d", follow.FollowUserId))
	}

	return following, nil
}

// 从 MySQL 获取粉丝列表
func getFollowersFromMySQL(userID int64) ([]string, error) {
	db := utils.GetMysqlDB()

	var follows []models.Follow
	err := db.Where("follow_user_id = ?", userID).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	var followers []string
	for _, follow := range follows {
		followers = append(followers, fmt.Sprintf("%d", follow.UserId))
	}

	return followers, nil
}

// 从 MySQL 获取朋友列表（关注和粉丝的交集）
func getFriendsFromMySQL(userID int64) ([]string, error) {
	db := utils.GetMysqlDB()

	// 获取关注列表和粉丝列表的交集（朋友列表）
	var friends []string
	subQuery := db.Table("follows").Select("follow_user_id as id").Where("user_id = ?", userID)
	err := db.Table("follows").Where("follow_user_id = ? AND user_id IN (?)", userID, subQuery).Find(&friends).Error
	if err != nil {
		return nil, err
	}

	var friendIDs []string
	for _, friend := range friends {
		friendIDs = append(friendIDs, fmt.Sprintf("%d", friend))
	}

	return friendIDs, nil
}

// 将关注列表添加到 Redis
func addFollowingToRedis(userID int64, following []string) error {
	client := utils.GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%d", config.FollowSetKey, userID)

	_, err := client.SAdd(ctx, key, following).Result()
	if err != nil {
		return err
	}

	return nil
}

// 将粉丝列表添加到 Redis
func addFollowersToRedis(userID int64, followers []string) error {
	client := utils.GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%d", config.FollowerSetKey, userID)

	_, err := client.SAdd(ctx, key, followers).Result()
	if err != nil {
		return err
	}

	return nil
}

// 将朋友列表添加到 Redis
func addFriendsToRedis(userID int64, friends []string) error {
	client := utils.GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%d", config.FriendSetKey, userID)

	_, err := client.SAdd(ctx, key, friends).Result()
	if err != nil {
		return err
	}

	return nil
}
