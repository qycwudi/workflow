package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gogogo/internal/config"
	"gogogo/internal/model"
)

type ServiceContext struct {
	Config        config.Config
	GogogoKvModel model.GogogoKvModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlDataSource)
	return &ServiceContext{
		Config:        c,
		GogogoKvModel: model.NewGogogoKvModel(conn),
	}
}
