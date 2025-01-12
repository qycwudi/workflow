package dispatch

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/dispatch/broadcast"
	"workflow/internal/svc"
)

// 初始化订阅
func InitPubSub(ctx *svc.ServiceContext) {
	// 创建带取消的上下文
	subCtx, cancel := context.WithCancel(context.Background())

	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 开启协程处理数据源客户端同步事件
	go func() {
		sync := broadcast.NewDatasourceClientSync(ctx.DatasourceModel)
		if err := sync.Subscribe(subCtx, sync.Handler); err != nil {
			// 记录错误但不影响主程序
			logx.Errorf("SubscribeDatasourceClientSyncEvent failed: %s", err.Error())
			return
		}
	}()

	// 开启协程处理API策略同步事件
	go func() {
		sync := broadcast.NewApiLoadSync()
		if err := sync.Subscribe(subCtx, sync.Handler); err != nil {
			// 记录错误但不影响主程序
			logx.Errorf("SubscribeApiLoadSyncEvent failed: %s", err.Error())
			return
		}
	}()

	// 开启协程处理Job加载同步事件
	go func() {
		sync := broadcast.NewJobLoadSync()
		if err := sync.Subscribe(subCtx, sync.Handler); err != nil {
			// 记录错误但不影响主程序
			logx.Errorf("SubscribeJobLoadSyncEvent failed: %s", err.Error())
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
