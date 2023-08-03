package test

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/mq"
	"github.com/streadway/amqp"
	"log"
	"testing"
)

func TestRabbitMQ(t *testing.T) {

	//mq.RabbitMqUrl = "amqp://rabbitMqUser:SyjwljgR&d133@114.132.217.209:5672/"
	conn, err := amqp.Dial(mq.RabbitMqUrl)
	fmt.Println(mq.RabbitMqUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := "commentMQ"
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化
		false,     // 是否自动删除
		false,     // 是否具有排他性
		false,     // 是否阻塞等待
		nil,       // 额外属性
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	message := "Hello, RabbitMQ!"
	err = ch.Publish(
		"",     // 交换机名称
		q.Name, // 队列名称
		false,  // 如果设置为 true，则根据 routingKey 在队列中查找对应的队列名称
		false,  // 是否阻塞等待
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf("Sent message to queue %s: %s", queueName, message)

}

func TestConsumeRabbitMQ(t *testing.T) {
	conn, err := amqp.Dial("amqp://rabbitMqUser:SyjwljgR&d133@114.132.217.209:5672/") // 连接到 RabbitMQ 服务器
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel() // 打开一个通道
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := "commentMQ" // 要监听的队列名称
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化队列
		false,     // 非自动删除队列
		false,     // 非独占队列
		false,     // 不等待服务器响应
		nil,       // 额外的属性
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 接收消息
	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者标识，为空表示使用默认标识
		true,   // 自动应答，即消费消息后自动向 RabbitMQ 确认消息已接收
		false,  // 非独占队列
		false,  // 不等待服务器响应
		false,  // 额外的属性
		nil,    // 额外的属性
	)
	if err != nil {
		log.Fatalf("Failed to consume messages from queue: %v", err)
	}

	// 处理接收到的消息
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
	}

}

func TestPublishMQ(t *testing.T) {
	mq.LikeRMQ.Publish("hello world")
}
