package mq

import (
	"encoding/json"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
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
}

func (followMQ *FollowMQ) consumer(message <-chan amqp.Delivery) {
	for d := range message {
		// Handle the received message
		var data models.FollowMQToUser
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}
		// Now, process the follow action based on the data received
		follow := models.Follow{
			UserId:       data.UserId,
			FollowUserId: data.FollowUserId,
		}
		switch data.ActionType {
		case 1: // Follow action
			err := follow.Insert(utils.GetMysqlDB())
			if err != nil {
				log.Printf("Error inserting follow record: %v", err)
				continue
			}
		case 2: // Unfollow action
			err := follow.Delete(utils.GetMysqlDB())
			if err != nil {
				log.Printf("Error deleting follow record: %v", err)
				continue
			}
		default:
			log.Printf("Invalid action type received: %d", data.ActionType)
		}
	}
}

var followRMQ *FollowMQ

// InitFollowRabbitMQ 初始化rabbitMQ连接。
func InitFollowRabbitMQ() {
	followRMQ = NewFollowRabbitMQ()
	go followRMQ.Consumer()
}
