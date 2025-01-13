package event

import "context"

// DatasourceEventPublisher 定义数据源事件发布接口
type DatasourceEventPublisher interface {
	Publish(ctx context.Context) error
}

// DatasourceEventType 数据源事件类型
const (
	DatasourceClientSyncEvent = "event_datasource_client_sync"
)

var (
	// DefaultDatasourcePublisher 默认的数据源事件发布器
	DefaultDatasourcePublisher DatasourceEventPublisher
)

// PublishDatasourceSync 发布数据源同步事件
func PublishDatasourceSync(ctx context.Context) error {
	if DefaultDatasourcePublisher != nil {
		return DefaultDatasourcePublisher.Publish(ctx)
	}
	return nil
}
