package job

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/dispatch/job"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type JobOnOffLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobOnOffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobOnOffLogic {
	return &JobOnOffLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobOnOffLogic) JobOnOff(req *types.JobOnOffRequest) (resp *types.JobOnOffResponse, err error) {
	// 查询job
	jobDetail, err := l.svcCtx.JobModel.FindByJobId(l.ctx, req.JobId)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "查询job失败")
	}
	// 更新job状态
	switch req.Status {
	case model.JobStatusOn:
		jobDetail.Status = model.JobStatusOn
		jobDetail.UpdateTime = time.Now()
		// 启动job
		jobInstance := &job.ChainJob{JobId: jobDetail.JobId, CanvasId: jobDetail.WorkspaceId}
		err = job.DispatcherManager.AddJob(jobDetail.JobId, jobDetail.JobCron, jobInstance)
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "启动job失败")
		}
	case model.JobStatusOff:
		jobDetail.Status = model.JobStatusOff
		jobDetail.UpdateTime = time.Now()
		// 取消job
		job.DispatcherManager.RemoveJob(jobDetail.JobId)
	}
	err = l.svcCtx.JobModel.Update(l.ctx, jobDetail)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "更新job状态失败")
	}
	return &types.JobOnOffResponse{
		JobId:  jobDetail.JobId,
		Status: jobDetail.Status,
	}, nil
}
