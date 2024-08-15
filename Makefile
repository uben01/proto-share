build:
	docker build . -t proto-share

run:
	docker run -it --rm -v./sample-project:/app -w/app proto-share /bin/sh -c "proto-share -config=./proto-share.config.yml"