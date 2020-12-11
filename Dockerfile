FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.io,direct
WORKDIR /build
ENV CC=gcc
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache gcc musl-dev
RUN go build -o qasystem .
######################################################################
FROM debian:stretch-slim
COPY ./conf /conf
COPY ./wait-for.sh /

COPY --from=builder /build/qasystem /

RUN sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list
RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		netcat; \
        chmod 755 wait-for.sh

# CMD "QAsystem" "-init"