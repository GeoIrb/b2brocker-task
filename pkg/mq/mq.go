package mq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

// PublishFunction for publish message to mq.
type PublishFunction func(ctx context.Context, data []byte) error

// Handler message from mq.
type Handler func(ctx context.Context, data []byte)

// Consumer mq.
type Consumer map[string]consume

type consume struct {
	delivery <-chan amqp.Delivery
	handler  Handler
}

// RabbitMQ message queue.
type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel

	queue Consumer
}

// Publisher creates queue with queueName and returns publish function
func (q *RabbitMQ) Publisher(queueName string) (publish PublishFunction, err error) {
	if _, err = q.channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return
	}
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

	c := Consumer{handler: handler}
	c.delivery, err = q.channel.Consume(queueName, "", true, false, false, false, nil)
	q.queue[queueName] = c
	return
}

// Listen starts listening message queue
func (q *RabbitMQ) Listen(ctx context.Context) {
	for _, consumer := range q.queue {
		go func(consumer Consumer) {
			for {
				select {
				case d, ok := <-consumer.delivery:
					if ok {
						consumer.handler(ctx, d.Body)
						continue
					}
					return
				case <-ctx.Done():
					return
				}
			}
		}(consumer)
	}
}

// Disconnect disconnect from message queue.
func (q *RabbitMQ) Disconnect() {
	q.channel.Close()
	q.connection.Close()
	return
}

// NewRabbitMQ ...
func NewRabbitMQ(url string) (q *RabbitMQ, err error) {
	q = &RabbitMQ{
		queue: make(map[string]Consumer),
	}
	if q.connection, err = amqp.Dial(url); err == nil {
		if q.channel, err = q.connection.Channel(); err == nil {
			err = q.channel.Qos(1, 0, false)
		}
	}
	return
}
