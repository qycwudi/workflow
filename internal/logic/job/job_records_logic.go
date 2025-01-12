package job

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobRecordsLogic {
	return &JobRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobRecordsLogic) JobRecords(req *types.JobRecordsRequest) (resp *types.JobRecordsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
