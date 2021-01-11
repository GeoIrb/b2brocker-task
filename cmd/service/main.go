package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"b2broker-task/pkg/mq"
	"b2broker-task/pkg/service"
	"b2broker-task/pkg/service/mqhandler"
)

const serviceName = "service"

type configuration struct {
	MQurl            string `envconfig:"MQ_URL" default:"amqp://guest:guest@localhost:5672/"`
	MQQueueToService string `envconfig:"MQ_QUEUE_TO_SERVICE" default:"to-service"`
	MQQueueToProxy   string `envconfig:"MQ_QUEUE_TO_PROXY" default:"to-proxy"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	level.Info(logger).Log("initialization", serviceName)

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		level.Error(logger).Log("msg", "read configuration", "err", err)
		os.Exit(1)
	}

	rabbitMQ, err := mq.NewRabbitMQ(cfg.MQurl)
	if err != nil {
		level.Error(logger).Log("msg", "connect to Rabbit MQ", "url", cfg.MQurl, "err", err)
		os.Exit(1)
	}
	toServicePublish, err := rabbitMQ.Publisher(cfg.MQQueueToService)
	if err != nil {
		level.Error(logger).Log("msg", "Rabbit MQ publisher", "queue", cfg.MQQueueToProxy, "err", err)
		os.Exit(1)
	}
	svc := service.New()

	handler := mqhandler.NewHandlerServer(
		svc,
		mqhandler.NewHandlerTransport(),
		toServicePublish,
		logger,
	)
	if err := rabbitMQ.Consumer(cfg.MQQueueToService, handler); err != nil {
		level.Error(logger).Log("msg", "add consumer", "queue", cfg.MQQueueToService, "err", err)
		os.Exit(1)
	}

	go func() {
		level.Info(logger).Log("msg", "mq server turn on")
		rabbitMQ.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
