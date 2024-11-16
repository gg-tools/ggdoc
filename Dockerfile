FROM golang:1.22 AS builder

# Golang 镜像加速
ENV CGO_ENABLED=0
ENV GOPRIVATE=""
ENV GOPROXY="https://goproxy.cn,direct"
ENV GOSUMDB="sum.golang.google.cn"

WORKDIR /root

ADD . .
RUN [ ! -d "vendor" ] && go mod download all || echo "go mod download skipped..."
RUN go build -o main main.go

FROM alpine

# Alpine 系统镜像加速
RUN sed -e 's/dl-cdn[.]alpinelinux.org/mirrors.aliyun.com/g' -i~ /etc/apk/repositories
RUN apk add --update --no-cache busybox-extras

# 使用国内时区
ENV TZ Asia/Shanghai
RUN apk add tzdata alpine-conf && \
    /sbin/setup-timezone -z Asia/Shanghai && \
    apk del alpine-conf

WORKDIR /root

ADD docs .
ADD hacks/mime/mime.types /etc/mime.types
ADD hacks/mime/globs2 /usr/share/mime/globs2
COPY --from=builder /root/main apidoc-server
RUN chmod +x apidoc-server

ENTRYPOINT ["/root/apidoc-server"]
