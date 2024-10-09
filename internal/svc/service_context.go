package svc

import (
	"workflow/internal/config"
	"workflow/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                   config.Config
	GoGoGoKvModel            model.GogogoKvModel
	WorkSpaceModel           model.WorkspaceModel
	WorkSpaceTagModel        model.WorkspaceTagModel
	WorkspaceTagMappingModel model.WorkspaceTagMappingModel
	ModuleModel              model.ModuleModel
	CanvasModel              model.CanvasModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlDataSource)
	return &ServiceContext{
		Config:                   c,
		GoGoGoKvModel:            model.NewGogogoKvModel(conn),
		WorkSpaceModel:           model.NewWorkspaceModel(conn),
		WorkSpaceTagModel:        model.NewWorkspaceTagModel(conn),
		WorkspaceTagMappingModel: model.NewWorkspaceTagMappingModel(conn),
		ModuleModel:              model.NewModuleModel(conn),
		CanvasModel:              model.NewCanvasModel(conn),
	}
}
