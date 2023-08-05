package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Configuration struct {
	MySQL       string            `yaml:"MySQL"`
	VideoServer VideoServerConfig `yaml:"VideoServer"`
	Redis       RedisConfig       `yaml:"Redis"`
	RabbitMQ    RabbitMQConfig    `yaml:"RabbitMQ"`
}

type RedisConfig struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
	PoolSize int    `yaml:"PoolSize"`
}

type RabbitMQConfig struct {
	Addr     string `yaml:"Addr"`
	User     string `yaml:"User"`
	Port     string `yaml:"Port"`
	Password string `yaml:"Password"`
}

type VideoServerConfig struct {
	Addr2 string `yaml:"Addr2"` //拼接play_url
	Addr  string `yaml:"Addr"`
	Api   struct {
		Upload struct {
			Path   string `yaml:"Path"`   // /ftpServer/upload/
			Method string `yaml:"Method"` // POST
		} `yaml:"Upload"`
	} `yaml:"Api"`
}

var Config Configuration

var (
	DefaultPage = "1"
	DefaultSize = "20"
	VideoCount  = 5
	BufferSize  = 1000
	// redis
	TokenTTL    float64 = 3600
	TokenKey    string  = "token:"
	LikeKey     string  = "Like:"
	LikeKeyTTL  float64 = 3600
	LikeLock    string  = "likeLock"
	UnLikeLock  string  = "unLikeLock"
	LikeLockTTL float64 = 60
	UserKey     string  = "user:"
	UsedrKeyTTL float64 = 3600

	//filter
	WordDictPath = "./public/sensitiveDict.txt"
)

var MailPassword = os.Getenv("MailPassword")

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

// 首字母大写其他的包才能调用
func ReadConfig() {
	configFile, err := ioutil.ReadFile("config/configuration.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = yaml.Unmarshal(configFile, &Config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
