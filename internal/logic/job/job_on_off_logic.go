package job

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/dispatch/broadcast"
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
		err = broadcast.NewJobLoadSync().Publish(l.ctx, &broadcast.JobLoadSyncMsg{
			JobId:       req.JobId,
			RuleChain:   string(jobDetail.Dsl),
			JobCron:     jobDetail.JobCron,
			WorkspaceId: jobDetail.WorkspaceId,
			Type:        broadcast.JobLoadSyncTypeAdd,
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "发送加载链服务消息失败")
		}
	case model.JobStatusOff:
		err = broadcast.NewJobLoadSync().Publish(l.ctx, &broadcast.JobLoadSyncMsg{
			JobId: req.JobId,
			Type:  broadcast.JobLoadSyncTypeRemove,
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "发送取消链服务消息失败")
		}
	}
	// 更新job状态
	jobDetail.Status = req.Status
	jobDetail.UpdateTime = time.Now()
	err = l.svcCtx.JobModel.Update(l.ctx, jobDetail)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "更新job状态失败")
	}
	return &types.JobOnOffResponse{
		JobId:  jobDetail.JobId,
		Status: jobDetail.Status,
	}, nil
}
