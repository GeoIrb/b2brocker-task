package mq

import (
	"github.com/streadway/amqp"
)

// RabbitMQ message queue
type RabbitMQ struct {
	connect *amqp.Connection
}

// Disconnect disconnect from message queue
func (q *RabbitMQ) Disconnect() error {
	return q.connect.Close()
}

// NewRabbitMQ ...
func NewRabbitMQ(url string) (q *RabbitMQ, err error) {
	q = &RabbitMQ{}
	if q.connect, err = amqp.Dial(url); err != nil {
		return
	}
	return
}
