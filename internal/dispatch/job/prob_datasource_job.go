package job

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	asynq2 "workflow/internal/asynq"
	"workflow/internal/asynq/processor"
)

type ProbDatasourceJob struct {
}

const (
	ProbDatasourceJobName = "ProbDatasourceJob"
)

func (p *ProbDatasourceJob) Run() {
	info, err := asynq2.AsynqClient.Enqueue(asynq.NewTask(processor.TOPIC_DATA_SOURCE_CLIENT_PROBE, nil))
	if err != nil {
		logx.Errorf("Failed to enqueue prob datasource job, error: %v", err)
	}
	logx.Infof("Enqueued prob datasource job, info: %+v", info)
}
