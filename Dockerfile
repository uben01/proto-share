FROM alpine:latest AS builder

RUN apk add --no-cache \
    go

COPY ./src /compiler/src
COPY ./go.work /compiler/go.work

WORKDIR /compiler

RUN go build -o ./proto-share ./src/main.go

FROM alpine:latest AS release

COPY --from=builder /compiler/proto-share /bin/proto-share

RUN apk add --no-cache \
    protoc \
    nodejs \
    npm

RUN npm install -g protoc-gen-ts