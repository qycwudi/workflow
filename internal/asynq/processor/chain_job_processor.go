package processor

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/workflow"
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
	startTime := time.Now()
	logx.Infof("%s start at: %s", TOPIC_CHAIN_JOB, startTime.Format("2006-01-02 15:04:05"))

	var payload ChainJobPayload
	if err := sonic.Unmarshal(t.Payload(), &payload); err != nil {
		logx.Errorf("%s parse payload failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}

	// 执行任务链

	// 读取 metadata
	workspace, err := processor.workspaceModel.GetWorkspaceById(context.Background(), payload.CanvasId)
	if err != nil {
		logx.Errorf("%s get metadata failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}
	metadata := make(map[string]string)
	err = sonic.Unmarshal([]byte(workspace.Configuration), &metadata)
	if err != nil {
		logx.Errorf("%s parse metadata failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}
	// 读取参数
	job, err := processor.jobModel.FindOneByJobId(context.Background(), payload.JobId)
	if err != nil {
		logx.Errorf("%s get job failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}

	if job.Status == model.JobStatusOff {
		logx.Infof("%s job is off, jobId: %s", TOPIC_CHAIN_JOB, payload.JobId)
		return nil
	}

	params := make(map[string]any)
	err = sonic.Unmarshal([]byte(job.Params), &params)
	if err != nil {
		logx.Errorf("%s parse params failed: %v", TOPIC_CHAIN_JOB, err)
		return err
	}
	// 运行
	serialId, result, err := workflow.Run(ctx, t.ResultWriter().TaskID(), payload.CanvasId, params)
	if err != nil {
		logx.Errorf("chain run failed: %v", err)
		return err
	}
	logx.Infof("chain run result:%+v, serialId: %s", result, serialId)
	resultJson, err := sonic.Marshal(result)
	if err != nil {
		logx.Errorf("marshal result failed: %v", err)
		return err
	}
	// 保存记录
	jobRecord := model.JobRecord{
		JobId:    payload.JobId,
		JobName:  job.JobName,
		Status:   "success",
		TraceId:  serialId,
		Param:    job.Params,
		Result:   string(resultJson),
		ExecTime: time.Now(),
		Duration: int64(time.Since(startTime).Milliseconds()),
	}
	_, err = processor.jobRecordModel.Insert(context.Background(), &jobRecord)
	if err != nil {
		logx.Errorf("save job record failed: %v", err)
	}

	_, _ = t.ResultWriter().Write([]byte(resultJson))
	logx.Infof("%s end at: %s, serialId: %s", TOPIC_CHAIN_JOB, time.Now().Format("2006-01-02 15:04:05"), serialId)
	return nil
}

func NewChainJobProcessor(workspaceModel model.WorkspaceModel, jobModel model.JobModel, jobRecordModel model.JobRecordModel) *ChainJobProcessor {
	return &ChainJobProcessor{
		workspaceModel: workspaceModel,
		jobModel:       jobModel,
		jobRecordModel: jobRecordModel,
	}
}
