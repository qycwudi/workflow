package pubsub

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/cache"
	"workflow/internal/rulego"
)

const (
	ApiLoadSyncEvent = "event_api_load_sync"
)

func PublishApiLoadSyncEvent(ctx context.Context, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return cache.Redis.Publish(ctx, ApiLoadSyncEvent, payloadBytes)
}

func SubscribeApiLoadSyncEvent(ctx context.Context, handler func(ctx context.Context, msg *redis.Message)) error {
	subscriber := cache.Redis.Subscribe(ctx, ApiLoadSyncEvent)
	defer subscriber.Close()

	ch := subscriber.Channel()
	logx.Infof("subscribe %s", ApiLoadSyncEvent)
	for {
		select {
		case <-ctx.Done():
			logx.Infof("%s context done: %s", ApiLoadSyncEvent, ctx.Err())
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				logx.Infof("%s channel closed", ApiLoadSyncEvent)
				return nil
			}
			handler(ctx, msg)
		}
	}
}

type ApiLoadSyncMsg struct {
	ApiId     string `json:"api_id"`
	RuleChain string `json:"rule_chain"`
}

func ApiLoadSyncHandler(ctx context.Context, msg *redis.Message) {
	// 读取 msg 消息
	var syncMsg ApiLoadSyncMsg
	err := json.Unmarshal([]byte(msg.Payload), &syncMsg)
	if err != nil {
		logx.Errorf("ApiLoadSyncHandler unmarshal msg failed: %s", err.Error())
		return
	}

	// 加载链服务
	err = rulego.RoleChain.LoadApiServiceChain(syncMsg.ApiId, []byte(syncMsg.RuleChain))
	if err != nil {
		logx.Errorf("ApiLoadSyncHandler load chain failed: %s", err.Error())
		return
	}
	logx.Infof("ApiLoadSyncHandler load chain success: %s", syncMsg.ApiId)
}
