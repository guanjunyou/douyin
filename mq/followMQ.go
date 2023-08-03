package mq

import (
	"github.com/streadway/amqp"
	"log"
)

type FollowMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewFollowRabbitMQ   获取followMQ的对应管道。
func NewFollowRabbitMQ() *FollowMQ {
	followMQ := &FollowMQ{
		RabbitMQ:  *Rmq,
		queueName: "followMQ",
	}
	ch, err := followMQ.conn.Channel()
	followMQ.channel = ch
	Rmq.failOnErr(err, "获取通道失败")
	return followMQ
}

// Publish 关注操作的发布配置。
func (followMQ *FollowMQ) Publish(message string) {

	_, err := followMQ.channel.QueueDeclare(
		followMQ.queueName,
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

	err1 := followMQ.channel.Publish(
		followMQ.exchange,
		followMQ.queueName,
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

// Consumer 关注关系的消费逻辑。
func (followMQ *FollowMQ) Consumer() {

	_, err := followMQ.channel.QueueDeclare(followMQ.queueName, true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	messages, err1 := followMQ.channel.Consume(
		followMQ.queueName,
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
	go followMQ.consumer(messages)
	//forever := make(chan bool)
	log.Println(messages)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	//<-forever

}
func (followMQ *FollowMQ) consumer(message <-chan amqp.Delivery) {
	for d := range message {
		log.Println(string(d.Body))
	}
}

var followRMQ *FollowMQ

// InitFollowRabbitMQ 初始化rabbitMQ连接。
func InitFollowRabbitMQ() {
	followRMQ = NewFollowRabbitMQ()
	followRMQ.Publish("hello word !")
	go followRMQ.Consumer()
}
