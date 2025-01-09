package processor

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/pubsub"
)

const (
	TOPIC_DATA_SOURCE_CLIENT_SYNC = "datasource_client_sync"
)

type DatasourceClientSyncProcessor struct {
}

type DatasourceClientSyncPayload struct {
}

func (processor *DatasourceClientSyncProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	logx.Infof("%s start at: %s", TOPIC_DATA_SOURCE_CLIENT_SYNC, time.Now().Format("2006-01-02 15:04:05"))

	// 发布数据源客户端同步事件 通知数据源客户端同步
	err := pubsub.PublishDatasourceClientSyncEvent(ctx)
	if err != nil {
		logx.Errorf("%s publish event failed: %s", TOPIC_DATA_SOURCE_CLIENT_SYNC, err.Error())
		return err
	}

	logx.Infof("%s end at: %s", TOPIC_DATA_SOURCE_CLIENT_SYNC, time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func NewDatasourceClientSyncProcessor() *DatasourceClientSyncProcessor {
	return &DatasourceClientSyncProcessor{}
}
