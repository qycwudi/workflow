package job

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/cache"
	"workflow/internal/dispatch/broadcast"
	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type JobEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobEditLogic {
	return &JobEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobEditLogic) JobEdit(req *types.JobEditRequest) (resp *types.JobEditResponse, err error) {
	// 查询当前任务
	jobDetail, err := l.svcCtx.JobModel.FindByJobId(l.ctx, req.JobId)
	if err != nil {
		return nil, err
	}
	jobDetail.JobCron = req.JobCron
	jobDetail.Params = req.JobParam
	jobDetail.JobName = req.JobName
	jobDetail.JobDesc = req.JobDesc
	err = l.svcCtx.JobModel.Update(l.ctx, jobDetail)
	if err != nil {
		return nil, err
	}

	// 发送加载链服务消息
	err = broadcast.NewJobLoadSync().Publish(l.ctx, &broadcast.JobLoadSyncMsg{
		JobId:       req.JobId,
		RuleChain:   string(jobDetail.Dsl),
		JobCron:     req.JobCron,
		WorkspaceId: jobDetail.WorkspaceId,
		Type:        broadcast.JobLoadSyncTypeEdit,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "发送加载链服务消息失败")
	}
	// 删除redis缓存
	err = cache.Redis.Del(l.ctx, fmt.Sprintf(cache.EnvRedisKey, req.JobId))
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除Job环境变量缓存失败")
	}
	return &types.JobEditResponse{JobId: req.JobId}, nil
}
