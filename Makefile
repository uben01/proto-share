.PHONY: build run test

IMAGE_NAME = proto-share

build:
	docker build . -t $(IMAGE_NAME)

run:
	docker run --rm -v ./samples/sample-project:/app $(IMAGE_NAME) "-config=./proto-share.config.yml"

test:
	go test ./...