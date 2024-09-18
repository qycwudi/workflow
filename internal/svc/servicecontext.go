package svc

import (
	"gogogo/internal/config"
	"gogogo/internal/model"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	GogogoKvModel   model.GogogoKvModel
	AsynqTaskClient *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlDataSource)
	return &ServiceContext{
		Config: c,
		GogogoKvModel:   model.NewGogogoKvModel(conn),
	}
}
