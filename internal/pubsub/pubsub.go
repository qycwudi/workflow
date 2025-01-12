package pubsub

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"workflow/internal/model"
	"workflow/internal/svc"
)

var datasourceModel model.DatasourceModel

// 初始化订阅
func NewPubSub(ctx *svc.ServiceContext) {
	// 创建带取消的上下文
	subCtx, cancel := context.WithCancel(context.Background())
	// 添加依赖
	datasourceModel = ctx.DatasourceModel

	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 开启协程处理数据源客户端同步事件
	go func() {
		if err := SubscribeDatasourceClientSyncEvent(subCtx, DatasourceClientSyncHandler); err != nil {
			// 记录错误但不影响主程序
			return
		}
	}()

	// 开启协程处理API策略同步事件
	go func() {
		if err := SubscribeApiLoadSyncEvent(subCtx, ApiLoadSyncHandler); err != nil {
			// 记录错误但不影响主程序
			return
		}
	}()

	// 开启协程处理Job加载同步事件
	go func() {
		if err := SubscribeJobLoadSyncEvent(subCtx, JobLoadSyncHandler); err != nil {
			// 记录错误但不影响主程序
			return
		}
	}()

	fmt.Println("pubsub init success")
	// 等待退出信号
	go func() {
		<-sigChan
		cancel() // 收到信号后取消上下文
	}()
}
