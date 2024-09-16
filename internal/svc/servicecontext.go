package svc

import (
	"github.com/hibiken/asynq"
	asynq2 "gogogo/internal/asynq"
	"gogogo/internal/config"
	"gogogo/internal/model"
)

type ServiceContext struct {
	Config          config.Config
	GogogoKvModel   model.GogogoKvModel
	AsynqTaskClient *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// conn := sqlx.NewMysql(c.MySqlDataSource)
	asynqClient := asynq2.NewAsynqClient(c.RedisConfig)
	go func() {
		asynq2.NewAsynqServer(c.RedisConfig)
	}()

	asynq2.AsynqTaskContext = asynq2.AsynqTask{AsynqTaskClient: asynqClient}
	return &ServiceContext{
		Config: c,
		// GogogoKvModel:   model.NewGogogoKvModel(conn),
		AsynqTaskClient: asynqClient,
	}
}
