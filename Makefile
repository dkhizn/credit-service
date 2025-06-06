IMAGE_NAME = credit-service
CONTAINER_NAME = credit_service_container

.PHONY: test lint build run stop clean

test:
	go test ./...

lint:
	golangci-lint run

build:
	docker build -t $(IMAGE_NAME) . 

run:
	docker run -d --name $(CONTAINER_NAME) -p 8080:8080 $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME) && docker rm $(CONTAINER_NAME)

clean: stop
	docker rmi $(IMAGE_NAME)
