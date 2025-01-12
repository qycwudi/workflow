package asynq

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/asynq/processor"
	"workflow/internal/svc"
)

// 创建Asynq服务
func InitAsynqServer(ctx *svc.ServiceContext) {
	server := asynq.NewServerFromRedisClient(ctx.RedisClient, asynq.Config{
		// IsFailure: func(err error) bool {
		// 	// If resource is not available, it's a non-failure error  false不记重试次数
		// 	return lo.Ternary(errors.Is(err, ErrResourceNotAvailable), false, true)
		// },
		// RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
		// 	if errors.Is(e, ErrResourceNotAvailable) {
		// 		return time.Duration(30) * time.Second
		// 	}
		// 	return asynq.DefaultRetryDelayFunc(n, e, t)
		// },
		Concurrency: 10,
	})
	// 注册处理器
	mux := asynq.NewServeMux()
	// 数据源客户端探测
	mux.Handle(processor.TOPIC_DATA_SOURCE_CLIENT_PROBE, processor.NewDatasourceClientProbeProcessor(ctx.DatasourceModel))
	// 任务JOB执行器
	mux.Handle(processor.TOPIC_CHAIN_JOB, processor.NewChainJobProcessor(ctx.WorkSpaceModel, ctx.JobModel, ctx.JobRecordModel))
	// 启动服务
	go func() {
		if err := server.Run(mux); err != nil {
			logx.Errorf("could not run server: %v", err)
		}
	}()
	fmt.Println("asynq server init success")
}
