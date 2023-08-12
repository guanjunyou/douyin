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
