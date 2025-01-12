package asynq

import (
	"github.com/hibiken/asynq"

	"workflow/internal/svc"
)

var AsynqClient *asynq.Client

// 创建Asynq客户端
func InitAsynqClient(ctx *svc.ServiceContext) {
	AsynqClient = asynq.NewClientFromRedisClient(ctx.RedisClient)
}
