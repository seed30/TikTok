package rabbitMQ

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/seed30/TikTok/dao"
	"github.com/seed30/TikTok/models"
	"github.com/streadway/amqp"
)

type CommentMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

var RmqCommentAdd *CommentMQ
var RmqCommentDel *CommentMQ

// 初始化 RabbitMQ 数据
func InitCommentRabbitMQ() {
	RmqCommentAdd = NewCommentRabbitMQ("comment_add")
	go RmqCommentAdd.Consumer()
	RmqCommentDel = NewCommentRabbitMQ("comment_del")
	go RmqCommentDel.Consumer()
}

func NewCommentRabbitMQ(queueName string) *CommentMQ {
	commentMQ := &CommentMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}
	cha, err := commentMQ.conn.Channel()
	commentMQ.channel = cha
	Rmq.failOnError(err, " 获取通道失败")
	return commentMQ
}

// 定义生产者
func (l *CommentMQ) Publish(message string) {
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
func (l *CommentMQ) Consumer() {
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
	case "comment_add":
		//点赞消费队列
		go l.consumerCommentAdd(messages)
	case "comment_del":
		//取消赞消费队列
		go l.consumerCommentDel(messages)

	}
	<-forever
}

func (l *CommentMQ) consumerCommentAdd(messages <-chan amqp.Delivery) {
	// 该操作是在接收到add 的请求后，像数据库发送请求
	// 且这种方式不会导致数据库在短时间内接收到大量的持久化信息，从而保护了数据库

	for d := range messages {

		data := fmt.Sprintf("%s", d.Body)

		log.Println("Comment data : ", data)
		var comment models.Comment
		json.Unmarshal([]byte(data), &comment)
		log.Println("Comment json : ", comment)
		err := dao.SaveComment(&comment)
		if err != nil {
			log.Println("添加评论失败 ！")
			return
		}

	}
}

func (l *CommentMQ) consumerCommentDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 将获取到的d进行解析，并获取得到 userId 和 videoId
		params := fmt.Sprintf("%s", d.Body)
		commentId, _ := strconv.ParseInt(params, 10, 64)

		// 操作数据库
		// 首先根据userId 和videoId 查询是否存在这两者的关系

		comment, err := dao.FindCommentByCommentId(commentId)
		if err != nil {
			log.Println(errors.New("consumerCommentDel Error : 数据库中不存在该条评论"))
		} else {
			// 说明存在
			// 判断cancel是否为 1 ， 否则操作失败
			err := dao.DeletComment(&comment)
			if err != nil {
				log.Println("将评论从评论表中删除失败")
				return
			}
		}
	}
}
