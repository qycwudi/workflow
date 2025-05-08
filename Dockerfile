FROM docker.m.daocloud.io/golang:1.24.3-alpine AS builder
# docker build --platform linux/amd64 -t workflow:b4 .
# docker login --username=qiangyuecheng registry.cn-hangzhou.aliyuncs.com
# Qycssg00
# docker tag 605517d6a654 registry.cn-hangzhou.aliyuncs.com/jenkins_construct_images/workflow:f6
# docker push registry.cn-hangzhou.aliyuncs.com/jenkins_construct_images/workflow:f7
# goctl kube deploy --name workflow-back --namespace workflow --port 8888 --o workflow-back-deploy.yaml

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
RUN GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o /app/workflow workflow.go


FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/workflow /app/workflow
COPY --from=builder /app/etc /app/etc

CMD ["./workflow", "-f", "etc/workflow-api.yaml"]
