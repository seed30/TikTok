package rabbitMQ

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/seed30/TikTok/dao"
	"github.com/seed30/TikTok/models"
	"github.com/streadway/amqp"
)

// 本包下创建关于点赞的 RabbitMq操作
type LikeMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

var RmqLikeAdd *LikeMQ
var RmqLikeDel *LikeMQ

// 初始化 RabbitMQ 连接
func InitLikeRabbitMQ() {
	RmqLikeAdd = NewLikeRabbitMQ("like_add")
	go RmqLikeAdd.Consumer()

	RmqLikeDel = NewLikeRabbitMQ("like_del")
	go RmqLikeDel.Consumer()
}

// NewLikeRabbitMQ 获取likeMQ的对应队列。
func NewLikeRabbitMQ(queueName string) *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}
	cha, err := likeMQ.conn.Channel()
	likeMQ.channel = cha
	Rmq.failOnError(err, "获取通道失败")
	return likeMQ
}

// 定义生产者
func (l *LikeMQ) Publish(message string) {
	// 先 申明
	_, err := l.channel.QueueDeclare(
		l.queueName,
		// 是否持久化
		false,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		//是否阻塞
		false,
		// 额外属性
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
		},
	)
	if err1 != nil {
		panic(err1)
	}
}

// Consumer like 关系的消费逻辑
func (l *LikeMQ) Consumer() {
	// 声明通道
	_, err := l.channel.QueueDeclare(
		l.queueName,
		// 是否持久化
		false,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		//是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}

	// 接收消息
	messages, err1 := l.channel.Consume(
		l.queueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		false,
		false,
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}
	//在golang实现RabbitMQ时，使用 "forever := make(chan bool)" 的原因是为了保证消费者进程不会退出。
	//通常情况下，消费者在处理完消息队列中的所有消息后就会退出。为了避免这种情况，我们需要在消费者进程中实现一个死循环，
	//以确保消费者在接收到新的消息时一直保持在线状态。
	//因此，"forever := make(chan bool)" 的作用是创建一个双向通道，并将其初始化为不阻塞状态，
	//从而保证消费者进程可以在接收到新消息时一直保持在线状态。
	//在代码结束处，我们可以通过 <-forever 命令来等待通道的关闭，从而确保程序的正常退出。
	forever := make(chan bool)
	switch l.queueName {
	case "like_add":
		//点赞消费队列
		go l.consumerLikeAdd(messages)
	case "like_del":
		//取消赞消费队列
		go l.consumerLikeDel(messages)

	}
	<-forever
}

func (l *LikeMQ) consumerLikeAdd(messages <-chan amqp.Delivery) {
	// 该操作是在接收到add 的请求后，像数据库发送请求
	// 且这种方式不会导致数据库在短时间内接收到大量的持久化信息，从而保护了数据库

	for d := range messages {
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		// 操作数据库
		// 首先根据userId 和videoId 查询是否存在这两者的关系

		like := models.Like{
			UserId:  userId,
			VideoId: videoId,
		}
		exist := dao.FindLike(&like)
		if exist {
			// 说明存在
			// 判断cancel是否为 2 ， 否则操作失败
			if like.Cancel == 2 {
				// 说明数据正确， 此时取消点赞，那么需要将数据库中的Cancel修改为1
				like.Cancel = 1
				err := dao.UpdateLike(&like)
				if err != nil {
					log.Println("consumerLikAdd Error : 更新数据库失败！")
				}
			} else {
				log.Println(errors.New("consumerLikeAdd Error : 用户已经点赞了，不可重复点赞"))
			}
		} else {
			// 说明数据不存在， 那么就需要向数据库中添加该条点赞数据
			like.Cancel = 1
			err := dao.CreateLike(&like)
			if err != nil {
				log.Println(errors.New("consumerLikeAd Error : 用户点赞信息添加失败"))
			}
		}
	}
}

func (l *LikeMQ) consumerLikeDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 将获取到的d进行解析，并获取得到 userId 和 videoId
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		videoId, _ := strconv.ParseInt(params[1], 10, 64)

		// 操作数据库
		// 首先根据userId 和videoId 查询是否存在这两者的关系

		like := models.Like{
			UserId:  userId,
			VideoId: videoId,
		}
		exist := dao.FindLike(&like)
		if exist {
			// 说明存在
			// 判断cancel是否为 1 ， 否则操作失败
			if like.Cancel == 1 {
				// 说明数据正确， 此时取消点赞，那么需要将数据库中的Cancel修改为2
				like.Cancel = 2
				err := dao.UpdateLike(&like)
				if err != nil {
					log.Println("consumerLikDel Error : 更新数据库失败！")
				}
			} else {
				log.Println(errors.New("consumerLikeDel Error : 用户取消点赞错误"))
			}
		} else {
			log.Println(errors.New("consumerLikeDel Error : 数据库中不存在该条数据"))
		}

	}
}
