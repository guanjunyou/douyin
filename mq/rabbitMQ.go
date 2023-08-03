package mq

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/streadway/amqp"
	"log"
)

var RabbitMqUrl string = "amqp://" + config.Config.RabbitMQ.User + ":" + config.Config.RabbitMQ.Password + "@" + config.Config.RabbitMQ.Addr +
	":" + config.Config.RabbitMQ.Port + "/"

type RabbitMQ struct {
	conn  *amqp.Connection
	mqurl string
}

var Rmq *RabbitMQ

// InitRabbitMQ 初始化RabbitMQ的连接和通道。
func InitRabbitMQ() {

	Rmq = &RabbitMQ{
		mqurl: "amqp://rabbitMqUser:SyjwljgR&d133@114.132.217.209:5672/",
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	Rmq.conn = dial
	if err != nil {
		log.Println("连接失败")
	}
}

// 关闭mq通道和mq的连接。
func (r *RabbitMQ) destroy() {
	r.conn.Close()
}

// 连接出错时，输出错误信息。
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s\n", err, message)
		panic(fmt.Sprintf("%s:%s\n", err, message))
	}
}
