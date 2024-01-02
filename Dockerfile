FROM golang:alpine as builder

WORKDIR /build
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="dalefengs@gmail.com"

WORKDIR /app

COPY --from=0 /app/server ./
COPY --from=0 /app/config.docker.yaml ./

EXPOSE 8818
ENTRYPOINT ./server -c config.docker.yaml
