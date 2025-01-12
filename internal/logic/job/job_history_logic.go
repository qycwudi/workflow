package job

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobHistoryLogic {
	return &JobHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobHistoryLogic) JobHistory(req *types.JobHistoryRequest) (resp *types.JobHistoryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
