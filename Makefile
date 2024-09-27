.PHONY: build run test

IMAGE_NAME = utobence/proto-share
IMAGE_TAG = 0.1

build:
	docker build . -t $(IMAGE_NAME):$(IMAGE_TAG) -t $(IMAGE_NAME):latest

publish:
	docker push $(IMAGE_NAME) -a

run:
	docker run --rm -v ./samples/sample-project:/app $(IMAGE_NAME):$(IMAGE_TAG) -config=./proto-share.config.yml -verbose

test:
	go test ./...