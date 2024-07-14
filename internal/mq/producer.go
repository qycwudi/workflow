package mq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func SendMessage(ctx context.Context, producer golang.Producer, marshal []byte, keys []string, tag string, topic string) error {
	msg := &golang.Message{
		Topic: topic,
		Body:  marshal,
	}
	// set keys and tag
	msg.SetKeys(keys...)
	msg.SetTag(tag)
	// send message in sync
	logx.WithContext(ctx).Infof("send message: %s", marshal)
	resp, err := producer.Send(ctx, msg)
	if len(resp) > 0 {
		logx.WithContext(ctx).Infof("send message result:%+v \n", resp[0])
	}
	if err != nil {
		logx.WithContext(ctx).Errorf("send message err:%+v", err)
		return err
	}
	return nil
}

func CreateLlmMqProducer(endPoint string) golang.Producer {
	os.Setenv("mq.consoleAppender.enabled", "true")
	os.Setenv("rocketmq.client.logLevel", "error")
	// os.Setenv("rocketmq.client.logRoot", "")
	golang.ResetLogger()
	producer, err := golang.NewProducer(&golang.Config{
		Endpoint:      endPoint,
		NameSpace:     "gogogo",
		ConsumerGroup: "gogogo-llm",
		Credentials: &credentials.SessionCredentials{
			AccessKey:    "",
			AccessSecret: "",
		},
	},
		golang.WithTopics(OCR_V1_TOPIC, LLM_TRAIT_EXTRACT_V1_TOPIC),
	)
	if err != nil {
		logx.Error("RockerMQ生产者创建失败")
		logx.Error(err)
		return nil
	}
	err = producer.Start()
	if err != nil {
		logx.Error("RockerMQ生产者启动失败")
		logx.Error(err)
		return nil
	}
	return producer
}

func CreateOcrMqProducer(endPoint string) golang.Producer {
	os.Setenv("mq.consoleAppender.enabled", "true")
	os.Setenv("rocketmq.client.logLevel", "error")
	// os.Setenv("rocketmq.client.logRoot", "")
	golang.ResetLogger()
	producer, err := golang.NewProducer(&golang.Config{
		Endpoint:      endPoint,
		NameSpace:     "gogogo",
		ConsumerGroup: "gogogo-llm",
		Credentials: &credentials.SessionCredentials{
			AccessKey:    "",
			AccessSecret: "",
		},
	},
		golang.WithTopics(OCR_V1_TOPIC, LLM_TRAIT_EXTRACT_V1_TOPIC),
	)
	if err != nil {
		logx.Error("RockerMQ生产者创建失败")
		logx.Error(err)
		return nil
	}
	err = producer.Start()
	if err != nil {
		logx.Error("RockerMQ生产者启动失败")
		logx.Error(err)
		return nil
	}
	return producer
}

func CloseProducer(producer golang.Producer) {
	err := producer.GracefulStop()
	if err != nil {
		fmt.Println("rocketMQ producer Shutdown error:", err.Error())
	}
}
