package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"
	"github.com/valyala/fasthttp"

	"b2broker-task/pkg/mq"
	"b2broker-task/pkg/proxy"
	"b2broker-task/pkg/proxy/httprouter"
	"b2broker-task/pkg/proxy/mqhandler"
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
	svc := proxy.New(toServicePublish)

	router := httprouter.New(svc, logger)

	server := &fasthttp.Server{
		Handler:          router.Handler,
		DisableKeepalive: true,
	}

	go func() {
		level.Info(logger).Log("msg", "http server turn on", "port", cfg.Port)
		if err := server.ListenAndServe(":" + cfg.Port); err != nil {
			level.Error(logger).Log("msg", "http server turn on", "err", err)
			os.Exit(1)
		}
	}()

	handler := mqhandler.NewHandlerServer(
		svc,
		mqhandler.NewTransport(),
		logger,
	)
	if err := rabbitMQ.Consumer(cfg.MQQueueToProxy, handler); err != nil {
		level.Error(logger).Log("msg", "add consumer", "queue", cfg.MQQueueToProxy, "err", err)
		os.Exit(1)
	}

	go func() {
		level.Info(logger).Log("msg", "mq server turn on")
		rabbitMQ.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

	if err := rabbitMQ.Shoutdown(); err != nil {
		level.Info(logger).Log("msg", "mq shoutdown", "err", err)
	}

	if err := server.Shutdown(); err != nil {
		level.Info(logger).Log("msg", "http server shoutdown", "err", err)
	}
}
