package svc

import (
	"github.com/hibiken/asynq"
	asynq2 "gogogo/internal/asynq"
	"gogogo/internal/config"
	"gogogo/internal/model"
	model2 "gogogo/internal/model/mongo"
)

type ServiceContext struct {
	Config          config.Config
	GogogoKvModel   model.GogogoKvModel
	AsynqTaskClient *asynq.Client
	MGHotDataModel  model2.HotDataModel
	MGColdDataModel model2.ColdDataModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// conn := sqlx.NewMysql(c.MySqlDataSource)
	asynqClient := asynq2.NewAsynqClient(c.RedisConfig)
	go func() {
		asynq2.NewAsynqServer(c.RedisConfig)
	}()

	// mongoDB
	mgHotDataModel := model2.NewHotDataModel(c.MongoDbUrl)
	mgColdDataModel := model2.NewColdDataModel(c.MongoDbUrl)
	asynq2.AsynqTaskContext = asynq2.AsynqTask{MGHotDataModel: mgHotDataModel, MGColdDataModel: mgColdDataModel, AsynqTaskClient: asynqClient}
	return &ServiceContext{
		Config: c,
		// GogogoKvModel:   model.NewGogogoKvModel(conn),
		AsynqTaskClient: asynqClient,
		MGHotDataModel:  mgHotDataModel,
		MGColdDataModel: mgColdDataModel,
	}
}
