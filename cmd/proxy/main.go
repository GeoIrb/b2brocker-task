package main

type configuration struct {
	Port string `envconfig:"PORT" default:"8080"`

	MQurl            string `envconfig:"MQ_URL" default:"8080"`
	MQQueueToService string `envconfig:"MQ_QUEUE_TO_SERVICE" default:"to-service"`
	MQQueueToProxy   string `envconfig:"MQ_QUEUE_TO_PROXY" default:"to-proxy"`
}

func main() {
}
