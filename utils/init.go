package utils

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
)

/**
数据库配置移到这里，因为config文件要放公共的变量，不要提交这个文件
*/
// MySQL数据库配置
const (
	userName = "root"
	password = "123456" // 更改成自己密码之后请不要提交自己的密码
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "douyin"
)

var MysqlDNS = strings.Join([]string{userName, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8mb4&parseTime=True"}, "")

var RedisConfig = &redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",  // no password set
	DB:       0,   // use default DB
	PoolSize: 100, // 连接池大小
}

var DB = Init()

var RDB = InitRedisDB()

func Init() *gorm.DB {
	db, err := gorm.Open(mysql.Open(MysqlDNS), &gorm.Config{})
	if err != nil {
		log.Println("gorm Init Error : ", err)
	}
	return db
}

func InitRedisDB() *redis.Client {
	return redis.NewClient(RedisConfig)
}
