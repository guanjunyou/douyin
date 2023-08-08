package utils

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"time"
)

func SaveTokenToRedis(username string, token string, expiration time.Duration) error {
	client := GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%v", config.TokenKey, username)
	err := client.Set(ctx, key, token, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetTokenFromRedis(username string) (string, error) {
	client := GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%v", config.TokenKey, username)
	token, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

// RefreshToken 刷新token有效期
func RefreshToken(username string, expiration time.Duration) error {
	client := GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%v", config.TokenKey, username)
	err := client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// FollowUser 关注某个用户
func FollowUser(followerID, followeeID int64) error {
	client := GetRedisDB()
	ctx := context.Background()

	// 将关注者ID添加到被关注者的关注集合中
	followeeKey := fmt.Sprintf("%v%d", config.FollowSetKey, followeeID)
	exists, err := client.Exists(ctx, followeeKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		// 如果被关注者的关注集合不存在，先创建一个空集合
		_, err := client.SAdd(ctx, followeeKey, "").Result()
		if err != nil {
			return err
		}
	}

	// 将被关注者ID添加到关注者的粉丝集合中
	followerKey := fmt.Sprintf("%v%d", config.FollowerSetKey, followerID)
	exists, err = client.Exists(ctx, followerKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		// 如果关注者的粉丝集合不存在，先创建一个空集合
		_, err := client.SAdd(ctx, followerKey, "").Result()
		if err != nil {
			return err
		}
	}

	// 执行关注操作
	_, err = client.SAdd(ctx, followeeKey, followerID).Result()
	if err != nil {
		return err
	}

	_, err = client.SAdd(ctx, followerKey, followeeID).Result()
	if err != nil {
		// 如果添加关注者的粉丝集合失败，需要回滚之前对被关注者的关注集合的修改
		client.SRem(ctx, followeeKey, followerID)
		return err
	}

	return nil
}

// UnfollowUser 取消关注某个用户
func UnfollowUser(followerID, followeeID int64) error {
	client := GetRedisDB()
	ctx := context.Background()

	// 将关注者ID从被关注者的关注集合中移除
	followeeKey := fmt.Sprintf("%v%d", config.FollowSetKey, followeeID)
	_, err := client.SRem(ctx, followeeKey, followerID).Result()
	if err != nil {
		return err
	}

	// 将被关注者ID从关注者的粉丝集合中移除
	followerKey := fmt.Sprintf("%v%d", config.FollowerSetKey, followerID)
	_, err = client.SRem(ctx, followerKey, followeeID).Result()
	if err != nil {
		// 如果移除关注者的粉丝集合失败，需要回滚之前对被关注者的关注集合的修改
		client.SAdd(ctx, followeeKey, followerID)
		return err
	}

	return nil
}

// GetUserFollowing 获取某个用户的关注列表
func GetUserFollowing(userID int64) ([]int64, error) {
	client := GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%d", config.FollowSetKey, userID)
	following, err := client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
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
	client := GetRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%d", config.FollowerSetKey, userID)
	followers, err := client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
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
	client := GetRedisDB()
	ctx := context.Background()

	// 获取关注列表和粉丝列表的键名
	followingKey := fmt.Sprintf("%v%d", config.FollowSetKey, userID)
	followersKey := fmt.Sprintf("%v%d", config.FollowerSetKey, userID)

	// 使用SINTER命令获取交集（朋友列表）
	friends, err := client.SInter(ctx, followingKey, followersKey).Result()
	if err != nil {
		return nil, err
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
