build:
	docker build . -t proto-share

run:
	docker run --rm -v./samples/sample-project:/app -w/app proto-share /bin/sh -c "proto-share -config=./proto-share.config.yml"

test:
	go test ./...