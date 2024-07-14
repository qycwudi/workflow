package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	asynq2 "gogogo/internal/asynq"
	"gogogo/internal/config"
	"gogogo/internal/model"
	model2 "gogogo/internal/model/mongo"
)

type ServiceContext struct {
	Config          config.Config
	GogogoKvModel   model.GogogoKvModel
	AsynqTaskClient *asynq.Client
	MGDataModel     model2.DataModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlDataSource)
	asynqClient := asynq2.NewAsynqClient(c.RedisConfig)
	go func() {
		asynq2.NewAsynqServer(c.RedisConfig)
	}()

	// mongoDB
	mgDataModel := model2.NewDataModel(c.MongoDbUrl)
	asynq2.AsynqTaskContext = asynq2.AsynqTask{MGDataModel: mgDataModel}
	return &ServiceContext{
		Config:          c,
		GogogoKvModel:   model.NewGogogoKvModel(conn),
		AsynqTaskClient: asynqClient,
		MGDataModel:     mgDataModel,
	}
}
