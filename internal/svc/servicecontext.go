package svc

import (
	"github.com/apache/rocketmq-clients/golang"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gogogo/internal/config"
	"gogogo/internal/model"
	"gogogo/internal/mq"
)

type ServiceContext struct {
	Config        config.Config
	GogogoKvModel model.GogogoKvModel
	OcrProducer   golang.Producer
	LlmProducer   golang.Producer
	ocrConsumer   golang.SimpleConsumer
	llmConsumer   golang.SimpleConsumer
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlDataSource)

	// 创建生产者
	ocrProducerMQ := mq.CreateOcrMqProducer(c.RocketMQNameServerAddress)
	llmProducerMQ := mq.CreateLlmMqProducer(c.RocketMQNameServerAddress)

	// 创建消费者
	ocrConsumer := mq.CreateOCRConsumer(c.RocketMQNameServerAddress)
	llmConsumer := mq.CreateLLMConsumer(c.RocketMQNameServerAddress)

	return &ServiceContext{
		Config:        c,
		GogogoKvModel: model.NewGogogoKvModel(conn),
		OcrProducer:   ocrProducerMQ,
		LlmProducer:   llmProducerMQ,
		ocrConsumer:   ocrConsumer,
		llmConsumer:   llmConsumer,
	}
}

func GetOCRConsumer(svc *ServiceContext) golang.SimpleConsumer {
	return svc.ocrConsumer
}
func GetLLMConsumer(svc *ServiceContext) golang.SimpleConsumer {
	return svc.llmConsumer
}
