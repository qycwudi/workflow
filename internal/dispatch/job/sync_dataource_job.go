package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/dispatch/broadcast"
)

type SyncDatasourceJob struct {
}

const (
	SyncDatasourceJobName = "DataSourceClientSync"
)

func (p *SyncDatasourceJob) Run() {
	err := broadcast.NewDatasourceClientSyncConstructor().Publish(context.Background())
	if err != nil {
		logx.Errorf("Failed to publish datasource client sync event, error: %v", err)
	}
}
