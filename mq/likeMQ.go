package mq

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/streadway/amqp"
)

type LikeMQ struct {
	RabbitMQ
	Channel        *amqp.Channel
	QueueUserName  string
	QueueVideoName string
	exchange       string
	key            string
}

// 初始化 channel
var LikeChannel chan models.LikeMQToVideo

func MakeLikeChannel() {
	ch := make(chan models.LikeMQToVideo, config.BufferSize)
	LikeChannel = ch
}

// NewLikeRabbitMQ 获取likeMQ的对应队列。
func NewLikeRabbitMQ() *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:       *Rmq,
		QueueUserName:  "userLikeMQ",
		QueueVideoName: "videoLikeMQ",
		exchange:       "likeExchange",
	}
	ch, err := likeMQ.conn.Channel()
	likeMQ.Channel = ch
	Rmq.failOnErr(err, "获取通道失败")
	return likeMQ
}

// Publish like操作的发布配置。
func (l *LikeMQ) Publish(message string) {
	//声明交换机
	err := l.Channel.ExchangeDeclare(
		//1.交换机名称
		l.exchange,
		//2、kind:交换机类型
		//	//amqp.ExchangeDirect 定向
		//	//amqp.ExchangeFanout 扇形（广播），发送消息到每个队列
		//	//amqp.ExchangeTopic 通配符的方式
		//	//amqp.ExchangeHeaders 参数匹配
		amqp.ExchangeFanout,
		//是否持久化
		true,
		//自动删除
		false,
		//内部使用
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	_, err = l.Channel.QueueDeclare(
		l.QueueUserName,
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
	_, err = l.Channel.QueueDeclare(
		l.QueueVideoName,
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
	//绑定队列和交换机
	err = l.Channel.QueueBind(l.QueueUserName, "", l.exchange, false, nil)
	if err != nil {
		panic(err)
	}
	err = l.Channel.QueueBind(l.QueueVideoName, "", l.exchange, false, nil)
	if err != nil {
		panic(err)
	}

	err1 := l.Channel.Publish(
		l.exchange,
		"",
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

//// Consumer like关系的消费逻辑。
//func (l *LikeMQ) Consumer() {
//
//	_, err := l.Channel.QueueDeclare(l.queueName, true, false, false, false, nil)
//
//	if err != nil {
//		panic(err)
//	}
//
//	//2、接收消息
//	messages, err1 := l.Channel.Consume(
//		l.queueName,
//		//用来区分多个消费者
//		"",
//		//是否自动应答
//		true,
//		//是否具有排他性
//		false,
//		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
//		false,
//		//消息队列是否阻塞
//		false,
//		nil,
//	)
//	if err1 != nil {
//		panic(err1)
//	}
//	go l.consumer(messages)
//	//forever := make(chan bool)
//	log.Println(messages)
//
//	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")
//
//	//<-forever
//
//}
//func (l *LikeMQ) consumer(message <-chan amqp.Delivery) {
//	for d := range message {
//		log.Println(string(d.Body))
//	}
//}

var LikeRMQ *LikeMQ

// InitLikeRabbitMQ 初始化rabbitMQ连接。
func InitLikeRabbitMQ() {
	LikeRMQ = NewLikeRabbitMQ()
	//LikeRMQ.Publish("hello word !")
	//go LikeRMQ.Consumer()
}
