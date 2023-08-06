package mq

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/streadway/amqp"
	"log"
)

type CommentMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

var CommentChannel chan models.CommentMQToVideo

func MakeCommentChannel() {
	ch := make(chan models.CommentMQToVideo, config.BufferSize)
	CommentChannel = ch
}

// NewCommentRabbitMQ  获取commentMQ的对应管道。
func NewCommentRabbitMQ() *CommentMQ {
	commentMQ := &CommentMQ{
		RabbitMQ:  *Rmq,
		queueName: "commentMQ",
	}
	ch, err := commentMQ.conn.Channel()
	commentMQ.channel = ch
	Rmq.failOnErr(err, "获取通道失败")
	return commentMQ
}

// Publish 评论操作的发布配置。
func (commentMQ *CommentMQ) Publish(message string) {

	_, err := commentMQ.channel.QueueDeclare(
		commentMQ.queueName,
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

	err1 := commentMQ.channel.Publish(
		commentMQ.exchange,
		commentMQ.queueName,
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

// Consumer 评论关系的消费逻辑。
func (commentMQ *CommentMQ) Consumer() {

	_, err := commentMQ.channel.QueueDeclare(commentMQ.queueName, true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	messages, err1 := commentMQ.channel.Consume(
		commentMQ.queueName,
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
	go commentMQ.consumer(messages)
	//forever := make(chan bool)
	log.Println(messages)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	//<-forever

}
func (commentMQ *CommentMQ) consumer(message <-chan amqp.Delivery) {
	for d := range message {
		log.Println(string(d.Body))
	}
}

var commentRMQ *CommentMQ

// InitCommentRabbitMQ  初始化rabbitMQ连接。
func InitCommentRabbitMQ() {
	commentRMQ = NewCommentRabbitMQ()
	commentRMQ.Publish("hello word !")
	go commentRMQ.Consumer()
}
