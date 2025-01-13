package processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/rulego"
)

const (
	TOPIC_CHAIN_JOB = "chain_job"
)

type ChainJobProcessor struct {
	workspaceModel model.WorkspaceModel
	jobModel       model.JobModel
	jobRecordModel model.JobRecordModel
}

type ChainJobPayload struct {
	JobId    string `json:"jobId"`
	CanvasId string `json:"canvasId"`
}

func (processor *ChainJobProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	traceId := uuid.New().String()
	startTime := time.Now()
	logx.Infof("%s start at: %s, traceId: %s", TOPIC_CHAIN_JOB, startTime.Format("2006-01-02 15:04:05"), traceId)

	var payload ChainJobPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		logx.Errorf("%s parse payload failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}

	// 执行任务链
	logx.Infof("%s execute chain job, jobId: %s, canvasId: %s, traceId: %s", TOPIC_CHAIN_JOB, payload.JobId, payload.CanvasId, traceId)

	// 读取 metadata
	workspace, err := processor.workspaceModel.GetWorkspaceById(context.Background(), payload.CanvasId)
	if err != nil {
		logx.Errorf("%s get metadata failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}
	metadata := make(map[string]string)
	err = json.Unmarshal([]byte(workspace.Configuration), &metadata)
	if err != nil {
		logx.Errorf("%s parse metadata failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}
	metadata["traceId"] = traceId
	metadata["startTime"] = time.Now().Format("2006-01-02 15:04:05")

	// 读取参数
	job, err := processor.jobModel.FindOneByJobId(context.Background(), payload.JobId)
	if err != nil {
		logx.Errorf("%s get job failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}

	if job.Status == model.JobStatusOff {
		logx.Infof("%s job is off, jobId: %s, traceId: %s", TOPIC_CHAIN_JOB, payload.JobId, traceId)
		return nil
	}

	// 运行
	result := rulego.RoleChain.Run(payload.JobId, metadata, job.Params)
	logx.Infof("chain run result:%+v, traceId: %s", result, traceId)
	// 保存记录
	jobRecord := model.JobRecord{
		JobId:    payload.JobId,
		JobName:  job.JobName,
		Status:   "success",
		TraceId:  traceId,
		Param:    job.Params,
		Result:   string(result.Data),
		ExecTime: time.Now(),
		Duration: int64(time.Since(startTime).Milliseconds()),
	}
	_, err = processor.jobRecordModel.Insert(context.Background(), &jobRecord)
	if err != nil {
		logx.Errorf("save job record failed: %v", err)
	}

	_, _ = t.ResultWriter().Write([]byte(result.Data))
	logx.Infof("%s end at: %s, traceId: %s", TOPIC_CHAIN_JOB, time.Now().Format("2006-01-02 15:04:05"), traceId)
	return nil
}

func NewChainJobProcessor(workspaceModel model.WorkspaceModel, jobModel model.JobModel, jobRecordModel model.JobRecordModel) *ChainJobProcessor {
	return &ChainJobProcessor{
		workspaceModel: workspaceModel,
		jobModel:       jobModel,
		jobRecordModel: jobRecordModel,
	}
}
