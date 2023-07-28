package config

import (
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
)

var (
	DefaultPage = "1"
	DefaultSize = "20"
	VideoCount  = 0

	// redis
	TokenTTL float64 = 3600
	TokenKey string  = "token:"
)

var MailPassword = os.Getenv("MailPassword")

// var MysqlDNS = os.Getenv("MysqlDNS")

// MySQL数据库配置
const (
	userName = "root"
	password = "123456" // 更改成自己密码之后请不要提交自己的密码
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "douyin"
)

var RedisConfig = &redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",  // no password set
	DB:       0,   // use default DB
	PoolSize: 100, // 连接池大小
}

//var MysqlDNS = "root@tcp(127.0.0.1:3306)"

var MysqlDNS = strings.Join([]string{userName, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8mb4&parseTime=True"}, "")

type ProblemBasic struct {
	Identity          string      `json:"identity"`           // 问题表的唯一标识
	Title             string      `json:"title"`              // 问题标题
	Content           string      `json:"content"`            // 问题内容
	ProblemCategories []int       `json:"problem_categories"` // 关联问题分类表
	MaxRuntime        int         `json:"max_runtime"`        // 最大运行时长
	MaxMem            int         `json:"max_mem"`            // 最大运行内存
	TestCases         []*TestCase `json:"test_cases"`         // 关联测试用例表
}

type TestCase struct {
	Input  string `json:"input"`  // 输入
	Output string `json:"output"` // 输出
}

var (
	DateLayout            = "2006-01-02 15:04:05"
	ValidGolangPackageMap = map[string]struct{}{
		"bytes":   {},
		"fmt":     {},
		"math":    {},
		"sort":    {},
		"strings": {},
	}
)
