package pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/cache"
)

const (
	ApiTacticsSyncEvent = "event_api_tactics_sync"
)

func PublishApiTacticsSyncEvent(ctx context.Context, payload interface{}) error {
	return cache.Redis.Publish(ctx, ApiTacticsSyncEvent, payload)
}

func SubscribeApiTacticsSyncEvent(ctx context.Context, handler func(ctx context.Context, msg *redis.Message)) error {
	subscriber := cache.Redis.Subscribe(ctx, ApiTacticsSyncEvent)
	defer subscriber.Close()

	ch := subscriber.Channel()
	logx.Infof("subscribe %s", ApiTacticsSyncEvent)
	for {
		select {
		case <-ctx.Done():
			logx.Infof("%s context done: %s", ApiTacticsSyncEvent, ctx.Err())
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				logx.Infof("%s channel closed", ApiTacticsSyncEvent)
				return nil
			}
			handler(ctx, msg)
		}
	}
}

func ApiTacticsSyncHandler(ctx context.Context, msg *redis.Message) {

}
