FROM golang:1.16-alpine AS builder
ENV CGO_ENABLED=0
ENV GOPRIVATE=""
ENV GOPROXY="https://goproxy.cn,direct"
ENV GOSUMDB="sum.golang.google.cn"
WORKDIR /root/swagger-ui-server/

ADD . .
RUN go mod download \
    && go test --cover $(go list ./... | grep -v /vendor/) \
    && go build -o main main.go

FROM alpine
ENV TZ Asia/Shanghai
WORKDIR /root/

COPY --from=builder /root/swagger-ui-server/main main
COPY --from=builder /root/swagger-ui-server/statics/ ./statics
RUN chmod +x main

ENTRYPOINT ["/root/main"]
