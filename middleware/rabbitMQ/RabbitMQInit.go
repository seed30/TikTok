package rabbitMQ

import (
	"log"

	"github.com/seed30/TikTok/config"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn  *amqp.Connection
	mqUrl string
}

func (r *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (r *RabbitMQ) destroy() {
	r.conn.Close()
}

var Rmq *RabbitMQ

func InitRabbitMQ() {
	Rmq = &RabbitMQ{
		mqUrl: config.RabbitMQURL,
	}

	dial, err := amqp.Dial(Rmq.mqUrl)
	Rmq.failOnError(err, "连接RabbitMQ失败")
	Rmq.conn = dial

	InitLikeRabbitMQ()
	InitCommentRabbitMQ()
}
