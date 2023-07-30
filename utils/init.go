package utils

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// GetMysqlDB 需要使用数据库的时候直接创建一个连接 调用此方法即可/**
func GetMysqlDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.MySQL), &gorm.Config{})
	if err != nil {
		log.Println("gorm Init Error : ", err)
	}
	return db
}

// GetRedisDB 需要使用数据库的时候直接创建一个连接 调用此方法即可/**
func GetRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password, // no password set
		DB:       config.Config.Redis.DB,       // use default DB
		PoolSize: config.Config.Redis.PoolSize, // 连接池大小
	})
}
