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
