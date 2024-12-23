FROM m.daocloud.io/golang:1.23-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY ./etc /app/etc
# RUN go build -ldflags="-s -w" -o /app/workflow workflow.go
RUN GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o /app/workflow workflow.go



FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/workflow /app/workflow
COPY --from=builder /app/etc /app/etc

CMD ["./workflow", "-f", "etc/workflow-api.yaml"]
