package main

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"
	"github.com/valyala/fasthttp"

	"b2broker-task/pkg/mq"
	"b2broker-task/pkg/proxy"
	"b2broker-task/pkg/proxy/httprouter"
)

const serviceName = "proxy-service"

type configuration struct {
	Port string `envconfig:"PORT" default:"8080"`

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

	mq, err := mq.NewRabbitMQ(cfg.MQurl)
	if err != nil {
		level.Error(logger).Log("msg", "connect to Rabbit MQ", "url", cfg.MQurl, "err", err)
		os.Exit(1)
	}
	toServicePublish, err := mq.Publisher(cfg.MQQueueToService)
	if err != nil {
		level.Error(logger).Log("msg", "Rabbit MQ publisher", "queue", cfg.MQQueueToProxy, "err", err)
		os.Exit(1)
	}
	srv := proxy.New(toServicePublish)

	router := httprouter.New(srv, logger)

	server := &fasthttp.Server{
		Handler: router.Handler,
	}

	go func() {
		level.Info(logger).Log("msg", "start server", "port", cfg.Port)
		if err := server.ListenAndServe(":" + cfg.Port); err != nil {
			level.Error(logger).Log("server run failure error", err)
			os.Exit(1)
		}
	}()

	mqconsumer.

	mq.Consumer(cfg.MQQueueToProxy,)


}
