package asynq

import (
	"github.com/hibiken/asynq"

	"workflow/internal/svc"
)

// 创建Asynq客户端
func NewAsynqClient(ctx *svc.ServiceContext) *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: ctx.Config.Redis.Host, Password: ctx.Config.Redis.Password, DB: ctx.Config.Redis.DB})
	return client
}
