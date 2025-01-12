package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	rolego "workflow/internal/rulego"
)

const (
	TOPIC_CHAIN_JOB = "chain_job"
)

type ChainJobProcessor struct {
}

type ChainJobPayload struct {
	JobId    string `json:"jobId"`
	CanvasId string `json:"canvasId"`
	Params   string `json:"params"`
}

func (processor *ChainJobProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	logx.Infof("%s start at: %s", TOPIC_CHAIN_JOB, time.Now().Format("2006-01-02 15:04:05"))

	var payload ChainJobPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		logx.Errorf("%s parse payload failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}

	// 执行任务链
	logx.Infof("%s execute chain job, jobId: %s, canvasId: %s", TOPIC_CHAIN_JOB, payload.JobId, payload.CanvasId)
	// 读取 metadata

	// 读取参数

	// 创建运行记录
	metadata := map[string]string{
		"canvasId": payload.CanvasId,
	}

	// 运行
	result := rolego.RoleChain.Run(payload.JobId, metadata, payload.Params)
	logx.Infof("chain run result:%+v", result)

	var respData interface{}
	if err := json.Unmarshal([]byte(result.Data), &respData); err != nil {
		logx.Errorf("parsing result failure,err:%v", err)
	}
	// 保存记录

	logx.Infof("%s end at: %s", TOPIC_CHAIN_JOB, time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func NewChainJobProcessor() *ChainJobProcessor {
	return &ChainJobProcessor{}
}
