package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/event"
)

type SyncDatasourceJob struct {
}

const (
	SyncDatasourceJobName = "DataSourceClientSync"
)

func (p *SyncDatasourceJob) Run() {
	err := event.PublishDatasourceSync(context.Background())
	if err != nil {
		logx.Errorf("Failed to publish datasource client sync event, error: %v", err)
	}
}
