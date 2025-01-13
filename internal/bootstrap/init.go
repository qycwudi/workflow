package bootstrap

import (
	"workflow/internal/dispatch/broadcast"
	"workflow/internal/dispatch/job"
	"workflow/internal/svc"
)

// Initialize 初始化系统组件
func Initialize(ctx *svc.ServiceContext) {
	broadcast.InitPubSub(ctx)
	job.InitDcron(ctx)
}
