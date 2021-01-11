package mq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type consume struct {
	delivery <-chan amqp.Delivery
	handler  Handler
}

// PublishFunction for publish message to mq.
type PublishFunction func(ctx context.Context, data []byte) error

// Handler message from mq.
type Handler func(ctx context.Context, data []byte)

// RabbitMQ message queue.
type RabbitMQ struct {
	ctx    context.Context
	cancel context.CancelFunc

	connection *amqp.Connection
	channel    *amqp.Channel

	queue map[string]consume
}

// Publisher creates queue with queueName and returns publish function
func (q *RabbitMQ) Publisher(queueName string) (publish PublishFunction, err error) {
	publish = func(ctx context.Context, data []byte) error {
		return q.channel.Publish("", queueName, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})
	}
	return
}

// Consumer adds handler for queue
func (q *RabbitMQ) Consumer(queueName string, handler Handler) (err error) {
	if _, isExist := q.queue[queueName]; isExist {
		return fmt.Errorf("queue is exist")
	}

	c := consume{handler: handler}
	c.delivery, err = q.channel.Consume(queueName, "", true, false, false, false, nil)
	q.queue[queueName] = c
	return
}

// ListenAndServe starts listening message queue and serves them
func (q *RabbitMQ) ListenAndServe() {
	q.ctx, q.cancel = context.WithCancel(context.Background())
	for _, consumer := range q.queue {
		go func(consume consume) {
			for {
				select {
				case d := <-consume.delivery:
					consume.handler(q.ctx, d.Body)
				case <-q.ctx.Done():
					return
				}
			}
		}(consumer)
	}
}

// Shoutdown stops listenning and disconnect from message queue.
func (q *RabbitMQ) Shoutdown() error {
	q.cancel()
	if err := q.channel.Close(); err != nil {
		return err
	}
	return q.connection.Close()
}

// NewRabbitMQ ...
func NewRabbitMQ(url string) (q *RabbitMQ, err error) {
	q = &RabbitMQ{
		queue: make(map[string]consume),
	}
	if q.connection, err = amqp.Dial(url); err == nil {
		if q.channel, err = q.connection.Channel(); err == nil {
			err = q.channel.Qos(1, 0, false)
		}
	}
	return
}
