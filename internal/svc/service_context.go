package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"workflow/internal/config"
	"workflow/internal/model"
)

type ServiceContext struct {
	Config                   config.Config
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
	DatasourceModel          model.DatasourceModel
	LocksModel               model.LocksModel
	RedisClient              redis.UniversalClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlUrn)

	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{c.Redis.Host},
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	return &ServiceContext{
		Config:                   c,
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
		DatasourceModel:          model.NewDatasourceModel(conn),
		LocksModel:               model.NewLocksModel(conn),
		RedisClient:              redisClient,
	}
}
