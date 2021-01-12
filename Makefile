test:
	docker run --rm -d -p 15672:15672 -p 5672:5672 --name b2bbroker-task-mq rabbitmq:3-management
	go test ./...

lint:
	go fmt ./...
	go vet ./...
	golint -set_exit_status $(go list ./...)
	golangci-lint run -D errcheck
	go mod tidy

build-images:
	docker build -t b2broker-task-mq -f build/mq/Dockerfile .
	docker build -t b2broker-task-proxy -f build/proxy/Dockerfile .
	docker build -t b2broker-task-service -f build/service/Dockerfile .

run:
	docker run --rm -d -p 15672:15672 -p 5672:5672 --name b2broker-task-mq rabbitmq:3-management
	docker run --rm -d -e PORT=8080 -p 8080:8080 --name b2broker-task-proxy b2broker-task-proxy 
	docker run --rm -d -e --name b2broker-task-service b2broker-task-service
stop:
	docker stop b2broker-task-mq b2broker-task-proxy b2broker-task-service