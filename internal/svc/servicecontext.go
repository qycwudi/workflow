package svc

import (
	"gogogo/internal/config"
	"gogogo/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                   config.Config
	GoGoGoKvModel            model.GogogoKvModel
	WorkSpaceModel           model.WorkspaceModel
	WorkSpaceTagModel        model.WorkspaceTagModel
	WorkspaceTagMappingModel model.WorkspaceTagMappingModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlDataSource)
	return &ServiceContext{
		Config:                   c,
		GoGoGoKvModel:            model.NewGogogoKvModel(conn),
		WorkSpaceModel:           model.NewWorkspaceModel(conn),
		WorkSpaceTagModel:        model.NewWorkspaceTagModel(conn),
		WorkspaceTagMappingModel: model.NewWorkspaceTagMappingModel(conn),
	}
}
