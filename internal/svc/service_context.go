package svc

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"workflow/internal/config"
	"workflow/internal/middleware"
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
	RedisClient              redis.UniversalClient
	UsersModel               model.UsersModel
	RolesModel               model.RolesModel
	PermissionsModel         model.PermissionsModel
	UserRolesModel           model.UserRolesModel
	RolePermissionsModel     model.RolePermissionsModel
	PermissionMiddleware     func(next http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySqlUrn)

	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{c.Redis.Host},
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	permissionsModel := model.NewPermissionsModel(conn)

	return &ServiceContext{
		Config:               c,
		WorkSpaceModel:       model.NewWorkspaceModel(conn),
		WorkSpaceTagModel:    model.NewWorkspaceTagModel(conn),
		ModuleModel:          model.NewModuleModel(conn),
		CanvasModel:          model.NewCanvasModel(conn),
		ApiModel:             model.NewApiModel(conn),
		ApiRecordModel:       model.NewApiRecordModel(conn),
		ApiSecretKeyModel:    model.NewApiSecretKeyModel(conn),
		SpaceRecordModel:     model.NewSpaceRecordModel(conn),
		TraceModel:           model.NewTraceModel(conn),
		DatasourceModel:      model.NewDatasourceModel(conn),
		RedisClient:          redisClient,
		UsersModel:           model.NewUsersModel(conn),
		RolesModel:           model.NewRolesModel(conn),
		PermissionsModel:     permissionsModel,
		UserRolesModel:       model.NewUserRolesModel(conn),
		RolePermissionsModel: model.NewRolePermissionsModel(conn),
		PermissionMiddleware: middleware.NewPermissionMiddleware(permissionsModel).Handle,
	}
}
