package job

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	asynq2 "workflow/internal/asynq"
	"workflow/internal/asynq/processor"
)

type ChainJob struct {
	JobId    string `json:"jobId"`
	CanvasId string `json:"canvasId"`
}

const (
	ChainJobName = "ChainJob"
)

func (c *ChainJob) Run() {
	// submit task
	params := ChainJob{
		JobId:    c.JobId,
		CanvasId: c.CanvasId,
	}
	jsonBytes, _ := json.Marshal(params)

	info, err := asynq2.AsynqClient.Enqueue(asynq.NewTask(processor.TOPIC_CHAIN_JOB, jsonBytes, asynq.MaxRetry(4), asynq.Timeout(6*time.Hour)), asynq.Retention(24*time.Hour))
	if err != nil {
		logx.Errorf("Failed to submit task, taskId:%s, error:%v", c.JobId, err)
		return
	}
	logx.Infof("Successfully submitted task, taskId:%s, queue:%s, canvasId:%s", info.ID, info.Queue, c.CanvasId)
}
