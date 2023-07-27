package utils

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/go-redis/redis/v8"
	"time"
)

func SaveTokenToRedis(username string, token string, expiration time.Duration) error {
	client := InitRedisDB()
	ctx := context.Background()
	key := fmt.Sprintf("%v%v", config.TokenKey, username)

	err := client.Set(ctx, key, token, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetTokenFromRedis(client *redis.Client, username string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%v%v", config.TokenKey, username)

	token, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}
