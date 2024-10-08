FROM golang:alpine AS builder

WORKDIR /build
COPY . .

#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache git

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -ldflags "-X main.version=$(git describe --tags `git rev-list --tags --max-count=1`)" -o server .

FROM ubuntu:22.04

LABEL MAINTAINER="dalefengs@gmail.com"

WORKDIR /app

RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y tzdata \
    ca-certificates vim && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=builder /build/server ./
COPY --from=builder /build/config.docker.yaml ./

EXPOSE 8818

ENTRYPOINT ./server -c config.docker.yaml
