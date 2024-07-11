run-consumer:
	go run cmd/consumer/main.go

run-producer:
	go run cmd/producer/main.go

IMAGE_NAME=service-mesh-example
IMAGE_TAG=0.0.0
build:
	go mod tidy && \
	docker build -f infra/docker/Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .