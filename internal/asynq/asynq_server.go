package asynq

import (
	"github.com/hibiken/asynq"

	"workflow/internal/svc"
)

// 创建Asynq服务
func NewAsynqServer(ctx *svc.ServiceContext) *asynq.Server {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: ctx.Config.Redis.Host, Password: ctx.Config.Redis.Password, DB: ctx.Config.Redis.DB},
		asynq.Config{
			Concurrency: 20,
		},
	)

	return server
}
