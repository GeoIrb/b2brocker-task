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
	docker build -t proxy -f build/proxy/Dockerfile .
	docker build -t service -f build/service/Dockerfile .

run:
	docker run --rm -d -p 15672:15672 -p 5672:5672 --name b2broker-task-mq rabbitmq:3-management

stop:
	docker stop b2broker-task-mq