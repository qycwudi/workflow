package asynq

import (
	"github.com/hibiken/asynq"

	"workflow/internal/svc"
)

// 创建Asynq客户端
func NewAsynqClient(ctx *svc.ServiceContext) *asynq.Client {
	client := asynq.NewClientFromRedisClient(ctx.RedisClient)
	return client
}
