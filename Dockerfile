FROM alpine:3.20

RUN apk add --no-cache \
    protoc \
    go  \
    nodejs \
    npm

RUN npm install -g protoc-gen-ts

COPY ./src /compiler/src
COPY ./templates /compiler/templates
COPY ./go.work /compiler/go.work

WORKDIR /compiler

RUN go generate ./src/main.go
RUN go build -o /bin/proto-share ./src/main.go
