package broadcast

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/cache"
	"workflow/internal/workflow"
)

const (
	// ApiLoadSyncEvent API加载同步事件 当接口配置发生变化后,需要同步服务的数据源连接池
	ApiLoadSyncEvent = "event_api_load_sync"
)

type ApiLoadSyncMsg struct {
	ApiId     string `json:"api_id"`
	RuleChain string `json:"rule_chain"`
}

type ApiLoadSync struct{}

func NewApiLoadSync() *ApiLoadSync {
	return &ApiLoadSync{}
}

func (a *ApiLoadSync) Publish(ctx context.Context, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logx.Errorf("[ApiLoadSync] marshal payload failed: %v", err)
		return err
	}
	return cache.Redis.Publish(ctx, ApiLoadSyncEvent, payloadBytes)
}

func (a *ApiLoadSync) Subscribe(ctx context.Context, handler func(ctx context.Context, msg *redis.Message)) error {
	subscriber := cache.Redis.Subscribe(ctx, ApiLoadSyncEvent)
	defer subscriber.Close()

	ch := subscriber.Channel()
	logx.Infof("[ApiLoadSync] start subscribing event: %s", ApiLoadSyncEvent)
	for {
		select {
		case <-ctx.Done():
			logx.Infof("[ApiLoadSync] context done, event: %s, error: %v", ApiLoadSyncEvent, ctx.Err())
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				logx.Infof("[ApiLoadSync] channel closed, event: %s", ApiLoadSyncEvent)
				return nil
			}
			// Handler
			handler(ctx, msg)
		}
	}
}

func (a *ApiLoadSync) Handler(ctx context.Context, msg *redis.Message) {
	logx.Infof("[ApiLoadSync] receive message, payload: %s", msg.Payload)
	// 读取 msg 消息
	var syncMsg ApiLoadSyncMsg
	err := json.Unmarshal([]byte(msg.Payload), &syncMsg)
	if err != nil {
		logx.Errorf("[ApiLoadSync] unmarshal message failed: %v, payload: %s", err, msg.Payload)
		return
	}

	// 注册任务流
	err = workflow.Register(ctx, syncMsg.ApiId, syncMsg.RuleChain)
	if err != nil {
		logx.Errorf("[ApiLoadSync] register workflow failed: %v, apiId: %s", err, syncMsg.ApiId)
		return
	}

	logx.Infof("[ApiLoadSync] load chain success, apiId: %s", syncMsg.ApiId)
}
