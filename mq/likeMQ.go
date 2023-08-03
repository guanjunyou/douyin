package mq

import (
	"github.com/streadway/amqp"
	"log"
)

type LikeMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewLikeRabbitMQ 获取likeMQ的对应队列。
func NewLikeRabbitMQ() *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:  *Rmq,
		queueName: "likeMQ",
	}
	ch, err := likeMQ.conn.Channel()
	likeMQ.channel = ch
	Rmq.failOnErr(err, "获取通道失败")
	return likeMQ
}

// Publish like操作的发布配置。
func (l *LikeMQ) Publish(message string) {

	_, err := l.channel.QueueDeclare(
		l.queueName,
		//是否持久化
		true,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}

	err1 := l.channel.Publish(
		l.exchange,
		l.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		panic(err)
	}

}

// Consumer like关系的消费逻辑。
func (l *LikeMQ) Consumer() {

	_, err := l.channel.QueueDeclare(l.queueName, true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	messages, err1 := l.channel.Consume(
		l.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}
	go consumer(messages)
	//forever := make(chan bool)
	log.Println(messages)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	//<-forever

}
func consumer(message <-chan amqp.Delivery) {
	for d := range message {
		log.Println(string(d.Body))
	}
}

var LikeRMQ *LikeMQ

// InitLikeRabbitMQ 初始化rabbitMQ连接。
func InitLikeRabbitMQ() {
	LikeRMQ = NewLikeRabbitMQ()
	LikeRMQ.Publish("hello word !")
	go LikeRMQ.Consumer()
}
