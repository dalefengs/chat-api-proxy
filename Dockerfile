FROM golang:alpine as builder

WORKDIR /build
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM ubuntu:22.04

LABEL MAINTAINER="dalefengs@gmail.com"

WORKDIR /app

COPY --from=builder /build/server ./
COPY --from=builder /build/config.docker.yaml ./

EXPOSE 8818
ENTRYPOINT ./server -c config.docker.yaml
