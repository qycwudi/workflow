package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type JobListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobListLogic {
	return &JobListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobListLogic) JobList(req *types.JobPublishListRequest) (resp *types.JobPublishListResponse, err error) {
	// 查询job列表
	total, jobList, err := l.svcCtx.JobModel.FindPage(l.ctx, req.JobName, req.WorkspaceId, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "查询job列表失败")
	}
	jobListResp := make([]types.JobPublishList, len(jobList))
	for i, job := range jobList {

		jobListResp[i] = types.JobPublishList{
			JobId:      job.JobId,
			JobName:    job.JobName,
			JobDesc:    job.JobDesc,
			JobCron:    job.JobCron,
			JobParam:   job.Params,
			Status:     job.Status,
			CreateTime: job.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: job.UpdateTime.Format("2006-01-02 15:04:05"),
		}
	}
	return &types.JobPublishListResponse{
		Total:    total,
		Current:  req.Current,
		PageSize: req.PageSize,
		List:     jobListResp,
	}, nil
}
