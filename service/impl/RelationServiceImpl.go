package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/mq"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"gopkg.in/errgo.v2/errors"
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
	//分布式锁 不能让用户连续两次关注或者取消关注同一个用户的请求进入
	userIdStr := strconv.FormatInt(userId, 10)
	toUserIdStr := strconv.FormatInt(toUserId, 10)

	lockKey := config.FollowLock + userIdStr + toUserIdStr
	unFollowLockKey := config.UnFollowLock + userIdStr + toUserIdStr

	if actionType == 1 {
		isSuccess, _ := utils.GetRedisDB().SetNX(context.Background(), lockKey, "0", time.Duration(config.FollowLockTTL)*time.Second).Result()
		if isSuccess == false {
			log.Println("已关注")
			return errors.New("已关注")
		} else {
			utils.GetRedisDB().Del(context.Background(), unFollowLockKey)
		}
	} else {
		isSuccess, _ := utils.GetRedisDB().SetNX(context.Background(), unFollowLockKey, "0", time.Duration(config.FollowLockTTL)*time.Second).Result()
		if isSuccess == false {
			log.Println("已取消关注")
			return errors.New("已取消关注")
		} else {
			utils.GetRedisDB().Del(context.Background(), lockKey)
		}
	}

	var isExists bool = false
	var err error
	var follow *models.Follow

	userFollowKey := config.FollowKey + userIdStr
	//看看缓存中有没有这个集合
	exits, _ := utils.GetRedisDB().Exists(context.Background(), userFollowKey).Result()
	if exits != 0 {
		// 看看这个集合中有没有这个ID
		result, _ := utils.GetRedisDB().SIsMember(context.Background(), userFollowKey, userIdStr).Result()
		isExists = result
		// 如果缓存里面的 Set 里面没有就要从数据库里面查
		if !isExists {
			follow, err = getFollowByUserIdAndToUserId(userId, toUserId)
			if err != nil {
				log.Printf("查询关注记录发生异常 = %v", err)
				return err
			}
			if follow.Id != 0 {
				isExists = true
			}
		}
	} else {
		// 缓存中没有则从数据库找
		follow, err = getFollowByUserIdAndToUserId(userId, toUserId)
		if err != nil {
			log.Printf("查询关注记录发生异常 = %v", err)
			return err
		}

		if follow.Id != 0 {
			isExists = true
		}

	}

	if actionType == 1 {

		if isExists {
			log.Printf("该用户已关注")
			//tx.Rollback()
			err = errors.New("已关注")
			return err
		}

		//mqData := models.LikeMQToVideo{UserId: userId, VideoId: videoId, ActionType: actionType}
		mqData := models.FollowMQToUser{UserId: userId, FollowUserId: toUserId, ActionType: actionType}
		// 加入 channel
		mq.FollowChannel <- mqData
		jsonData, err := json.Marshal(mqData)
		if err != nil {
			log.Println("json序列化失败 = #{err}")
			//TODO 处理失败导致的数据不一致
		}
		//加入消息队列
		mq.FollowRMQ.Publish(string(jsonData))

		return nil

	} else if actionType == 2 {

		if !isExists && (follow == nil || follow.Id == 0) {
			log.Printf("未找到要取消的点赞记录")
			err = errors.New("-2")
			//tx.Rollback()
			return err
		}

		//mqData := models.LikeMQToVideo{UserId: userId, VideoId: videoId, ActionType: actionType}
		mqData := models.FollowMQToUser{UserId: userId, FollowUserId: toUserId, ActionType: actionType}
		// 加入 channel
		mq.FollowChannel <- mqData
		jsonData, err := json.Marshal(mqData)
		if err != nil {
			log.Printf("json序列化失败 = #{err}")
			//TODO 处理失败导致的数据不一致
		}
		//加入消息队列
		mq.FollowRMQ.Publish(string(jsonData))
		// TODO 消息队列处理失败会导致数据不一致

		return nil

	}
	return nil
}

func FollowConsumer(ch <-chan models.FollowMQToUser) {
	for {
		select {
		case msg := <-ch:
			// 在这里处理接收到的消息
			tx := utils.GetMysqlDB().Begin()
			if msg.ActionType == 1 {
				follow := models.Follow{
					CommonEntity: utils.NewCommonEntity(),
					UserId:       msg.UserId,
					FollowUserId: msg.FollowUserId,
				}
				err1 := follow.Insert(tx)
				if err1 != nil {
					log.Printf(err1.Error())
					tx.Rollback()
				}
				tx.Commit()

				userIdStr := strconv.FormatInt(msg.UserId, 10)
				toUserIdStr := strconv.FormatInt(msg.FollowUserId, 10)
				followSetKey := config.FollowKey + userIdStr
				followerSetKey := config.FollowerKey + toUserIdStr
				exists, _ := utils.GetRedisDB().Exists(context.Background(), followSetKey).Result()
				if exists == 0 { //缓存里面没有
					errBuildRedis := BuildFollowRedis(msg.UserId)
					if errBuildRedis != nil {
						log.Println("重建缓存失败", errBuildRedis)
					}
				} else {
					utils.GetRedisDB().SAdd(context.Background(), followSetKey, toUserIdStr)
				}

				exists1, _ := utils.GetRedisDB().Exists(context.Background(), followerSetKey).Result()
				if exists1 == 0 { //缓存里面没有
					errBuildRedis := BuildFollowerRedis(msg.FollowUserId)
					if errBuildRedis != nil {
						log.Println("重建缓存失败", errBuildRedis)
					}
				} else {
					utils.GetRedisDB().SAdd(context.Background(), followerSetKey, userIdStr)
				}
			}

			if msg.ActionType == 2 {
				follow, err := getFollowByUserIdAndToUserId(msg.UserId, msg.FollowUserId)
				if err != nil {
					tx.Rollback()
				}
				err1 := follow.Delete(tx)
				if err1 != nil {
					log.Printf(err1.Error())
					tx.Rollback()
				}
				tx.Commit()

				userIdStr := strconv.FormatInt(msg.UserId, 10)
				toUserIdStr := strconv.FormatInt(msg.FollowUserId, 10)
				followSetKey := config.FollowKey + userIdStr
				followerSetKey := config.FollowerKey + toUserIdStr
				// 删除缓存 SET 中的ID ， 避免脏数据产生
				utils.GetRedisDB().SRem(context.Background(), followSetKey, toUserIdStr)
				utils.GetRedisDB().SRem(context.Background(), followerSetKey, userIdStr)
			}
		default:
			// 如果channel为空，暂停一段时间后重新监听
			time.Sleep(time.Millisecond * 1)
		}
	}
}

// 重建缓存
func BuildFollowRedis(userId int64) error {
	relationService := RelationServiceImpl{}
	relationService.Logger = logrus.New()
	idSet, err := relationService.GetFollows(userId)
	if err != nil {
		return err
	}
	userIdStr := strconv.FormatInt(userId, 10)
	followSetKey := config.FollowKey + userIdStr
	var strValues []string
	for i := range idSet {
		strValues = append(strValues, strconv.FormatInt(idSet[i].Id, 10))
	}

	ctx := context.Background()
	err = utils.GetRedisDB().SAdd(ctx, followSetKey, strValues).Err()
	errSetTime := utils.GetRedisDB().Expire(ctx, followSetKey, time.Duration(config.FollowKeyTTL)*time.Second).Err()
	if errSetTime != nil {
		log.Println("redis 时间设置失败", errSetTime.Error())
	}
	return err
}

func BuildFollowerRedis(toUserId int64) error {
	relationService := RelationServiceImpl{}
	relationService.Logger = logrus.New()
	idSet, err := relationService.GetFollowers(toUserId)
	if err != nil {
		return err
	}
	toUserIdStr := strconv.FormatInt(toUserId, 10)
	followerSetKey := config.FollowerKey + toUserIdStr
	var strValues []string
	for i := range idSet {
		strValues = append(strValues, strconv.FormatInt(idSet[i].Id, 10))
	}

	ctx := context.Background()
	err = utils.GetRedisDB().SAdd(ctx, followerSetKey, strValues).Err()
	errSetTime := utils.GetRedisDB().Expire(ctx, followerSetKey, time.Duration(config.FollowKeyTTL)*time.Second).Err()
	if errSetTime != nil {
		log.Println("redis 时间设置失败", errSetTime.Error())
	}
	return err
}

// 创建消费者协程
func MakeFollowGroutine() {
	numConsumers := 20
	for i := 0; i < numConsumers; i++ {
		go FollowConsumer(mq.FollowChannel)
	}
}

// GetFollows 查询关注列表
func (relationServiceImpl RelationServiceImpl) GetFollows(userId int64) ([]models.User, error) {
	relationServiceImpl.Logger.Info("GetFollows\n")
	var users []models.User

	// 查询Redis中是否存在该用户ID，如果存在则，则使用协程获取数据并查询放到users中，否则原逻辑
	follows := utils.GetRedisDB().SMembers(context.Background(), fmt.Sprintf("follow:%d", userId)).Val()

	fmt.Println("follows = ", follows)

	if len(follows) > 0 {
		users = make([]models.User, len(follows))
		var wg sync.WaitGroup
		tx := utils.GetMysqlDB().Begin()
		for i, v := range follows {
			wg.Add(1)
			go func(index int, value string) {
				defer wg.Done()
				tx.Table("user").Where("id = ? and is_deleted = ?", value, 0).Find(users[index])
			}(i, v)
		}
		wg.Wait()
	} else {
		err := utils.GetMysqlDB().Table("follow").Where("user_id = ? AND is_deleted != ?", userId, 1).Find(&users).Error
		if err != nil {
			return nil, err
		}
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

	// 查询redis中有没有粉丝列表集合，如果有则启动多协程查询，如果没有则原逻辑
	followers := utils.GetRedisDB().SMembers(context.Background(), fmt.Sprintf("follower:%d", userId)).Val()
	fmt.Println("followers = ", followers)

	if len(followers) > 0 {
		users = make([]models.User, len(followers))
		var wg sync.WaitGroup
		tx := utils.GetMysqlDB().Begin()
		for i, v := range followers {
			wg.Add(1)
			go func(index int, value string) {
				defer wg.Done()
				tx.Table("user").Where("id = ? and is_deleted = ?", value, 0).Find(users[index])
			}(i, v)
		}
		wg.Wait()

		return users, nil
	}

	err := utils.GetMysqlDB().Table("follow").Where("follow_user_id = ? AND is_deleted != ?", userId, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetFriends 查询好友列表
func (relationServiceImpl RelationServiceImpl) GetFriends(userId int64) ([]models.User, error) {

	jugeExist(userId, "follow", fromMysqlToRedis)

	jugeExist(userId, "follower", fromMysqlToRedis)

	key1 := fmt.Sprintf("%s:%d", "follow", userId)
	key2 := fmt.Sprintf("%s:%d", "follower", userId)
	vals := utils.GetRedisDB().SInter(context.Background(), key1, key2).Val()

	users := make([]models.User, 0)
	err := utils.GetMysqlDB().Table("user").Find(&users, vals).Error

	return users, err

	// follows, err := relationServiceImpl.GetFollows(userId)
	// if err != nil {
	// 	return nil, err
	// }
	// followers, err := relationServiceImpl.GetFollowers(userId)
	// if err != nil {
	// 	return nil, err
	// }
	// var friends []models.User
	// for _, user := range followers {
	// 	if containsID(follows, user.Id) {
	// 		friends = append(friends, user)
	// 	}
	// }
	// return friends, nil
}

func fromMysqlToRedis(typeStr string, userId int64) (err error) {

	var wh string
	var s string

	if typeStr == "follow" {
		wh = "user_id = ? and is_deleted = ?"
		s = "follow_user_id"
	} else {
		wh = "follow_user_id = ? and is_deleted = ?"
		s = "user_id"
	}

	userIds := make([]int64, 0)
	if err = utils.GetMysqlDB().Table("follow").Select(s).Where(wh, userId, 0).Find(&userIds).Error; err != nil {
		return
	}

	key := fmt.Sprintf("%s:%d", typeStr, userId)
	if len(userIds) > 0 {
		redisDB := utils.GetRedisDB()
		redisDB.SAdd(context.Background(), key, userIds)
	}

	return
}

// 判断redis中是否存在key，并且不存在时，调用回调函数
func jugeExist(userId int64, typeStr string, callback func(t string, u int64) error) (err error) {
	followerExists := utils.GetRedisDB().Exists(context.Background(), fmt.Sprintf("%s:%d", typeStr, userId)).Val()

	if followerExists == 0 {
		// 不存在
		err = callback(typeStr, userId)
		return
	}
	return
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

func getFollowByUserIdAndToUserId(userId int64, toUserId int64) (*models.Follow, error) {
	res := &models.Follow{}
	err := utils.GetMysqlDB().Model(models.Follow{}).Where("user_id = ? AND follow_user_id = ? AND is_deleted = ?", userId, toUserId, 0).Find(res).Error
	return res, err
}
