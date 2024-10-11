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
	ApiModel                 model.ApiModel
	ApiRecordModel           model.ApiRecordModel
	ApiSecretKeyModel        model.ApiSecretKeyModel
	SpaceRecordModel         model.SpaceRecordModel
	TraceModel               model.TraceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlUrn)
	return &ServiceContext{
		Config:                   c,
		GoGoGoKvModel:            model.NewGogogoKvModel(conn),
		WorkSpaceModel:           model.NewWorkspaceModel(conn),
		WorkSpaceTagModel:        model.NewWorkspaceTagModel(conn),
		WorkspaceTagMappingModel: model.NewWorkspaceTagMappingModel(conn),
		ModuleModel:              model.NewModuleModel(conn),
		CanvasModel:              model.NewCanvasModel(conn),
		ApiModel:                 model.NewApiModel(conn),
		ApiRecordModel:           model.NewApiRecordModel(conn),
		ApiSecretKeyModel:        model.NewApiSecretKeyModel(conn),
		SpaceRecordModel:         model.NewSpaceRecordModel(conn),
		TraceModel:               model.NewTraceModel(conn),
	}
}
