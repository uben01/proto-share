FROM alpine:latest AS builder

RUN apk add --no-cache \
    go

COPY . /compiler/src

WORKDIR /compiler/src

RUN go build -o ../proto-share .

FROM alpine:latest AS release

COPY --from=builder /compiler/proto-share /bin/proto-share

RUN apk add --no-cache \
    protoc \
    nodejs \
    npm

RUN npm install -g protoc-gen-ts