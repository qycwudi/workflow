package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	asynq2 "gogogo/internal/asynq"
	"gogogo/internal/config"
	"gogogo/internal/model"
)

type ServiceContext struct {
	Config          config.Config
	GogogoKvModel   model.GogogoKvModel
	AsynqTaskClient *asynq.Client
	// rocketMQ太重,优先使用asynq https://github.com/hibiken/asynq
	// -OcrProducer   golang.Producer
	// -LlmProducer   golang.Producer
	// -ocrConsumer   golang.SimpleConsumer
	// -llmConsumer   golang.SimpleConsumer
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlDataSource)

	// -// 创建生产者
	// -ocrProducerMQ := mq.CreateOcrMqProducer(c.RocketMQNameServerAddress)
	// -llmProducerMQ := mq.CreateLlmMqProducer(c.RocketMQNameServerAddress)
	//
	// -// 创建消费者
	// -ocrConsumer := mq.CreateOCRConsumer(c.RocketMQNameServerAddress)
	// -llmConsumer := mq.CreateLLMConsumer(c.RocketMQNameServerAddress)

	asynqClient := asynq2.NewAsynqClient(c.RedisConfig)
	go func() {
		asynq2.NewAsynqServer(c.RedisConfig)
	}()
	return &ServiceContext{
		Config:          c,
		GogogoKvModel:   model.NewGogogoKvModel(conn),
		AsynqTaskClient: asynqClient,
		// -OcrProducer:   ocrProducerMQ,
		// -LlmProducer:   llmProducerMQ,
		// -ocrConsumer:   ocrConsumer,
		// -llmConsumer:   llmConsumer,
	}
}

// -func GetOCRConsumer(svc *ServiceContext) golang.SimpleConsumer {
// 	-return svc.ocrConsumer
// -}
// -func GetLLMConsumer(svc *ServiceContext) golang.SimpleConsumer {
// 	-return svc.llmConsumer
// -}
