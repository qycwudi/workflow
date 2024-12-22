package pubsub

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/cache"
	"workflow/internal/datasource"
	"workflow/internal/enum"
	"workflow/internal/model"
)

const (
	// 数据源客户端同步事件 当数据源配置发生变化后,需要同步服务的数据源连接池
	DatasourceClientSyncEvent = "event_datasource_client_sync"
)

func PublishDatasourceClientSyncEvent(ctx context.Context) error {
	return cache.Redis.Publish(ctx, DatasourceClientSyncEvent, "")
}

func SubscribeDatasourceClientSyncEvent(ctx context.Context, handler func(ctx context.Context, msg *redis.Message)) error {
	subscriber := cache.Redis.Subscribe(ctx, DatasourceClientSyncEvent)
	defer subscriber.Close()
	ch := subscriber.Channel()
	logx.Infof("subscribe %s", DatasourceClientSyncEvent)
	for {
		select {
		case <-ctx.Done():
			logx.Infof("%s context done: %s", DatasourceClientSyncEvent, ctx.Err())
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				logx.Infof("%s channel closed", DatasourceClientSyncEvent)
				return nil
			}
			handler(ctx, msg)
		}
	}
}

func DatasourceClientSyncHandler(ctx context.Context, msg *redis.Message) {
	logx.Infof("%s start at: %s", DatasourceClientSyncEvent, time.Now().Format("2006-01-02 15:04:05"))
	total, datasourceList, err := datasourceModel.FindDataSourcePageList(ctx, model.PageListBuilder{
		Switch: model.DatasourceSwitchOn,
	}, 1, 9999)
	if err != nil {
		logx.Error("Failed to get datasource list", err)
		return
	}

	logx.Infof("Successfully fetched %d datasources", total)
	successCount := 0
	failCount := 0
	skipCount := 0
	for _, ds := range datasourceList {
		// 跳过fileServer
		if ds.Type == enum.FileServerType.String() {
			skipCount++
			logx.Infof("%s skip: %d, %s", DatasourceClientSyncEvent, ds.Id, ds.Type)
			continue
		}
		err := datasource.DataSourcePool.UpdateDataSource(ds.Id, ds.Config, ds.Type, ds.Hash)
		if err != nil {
			logx.Errorf("%s DatasourceId: %d,DatasourceType: %s,DatasourceName: %s, update failed: %s", DatasourceClientSyncEvent, ds.Id, ds.Type, ds.Name, err.Error())
			failCount++
		} else {
			successCount++
		}
	}
	// 清理连接池中已删除的数据源连接
	poolIds := make(map[int64]bool)
	for _, ds := range datasourceList {
		poolIds[ds.Id] = true
	}

	// 遍历连接池,清理不存在于数据库的连接
	clearCount := 0
	for id := range datasource.DataSourcePool.GetHash() {
		if !poolIds[id] {
			if err := datasource.DataSourcePool.ClearDataSource(id); err != nil {
				logx.Errorf("clear datasource: %d, err: %v", id, err)
			} else {
				clearCount++
				logx.Infof("clear datasource: %d, success", id)
			}
		}
	}

	logx.Infof("%s end at: %s, success: %d, failed: %d, skip: %d",
		DatasourceClientSyncEvent, time.Now().Format("2006-01-02 15:04:05"), successCount, failCount, skipCount)
}
