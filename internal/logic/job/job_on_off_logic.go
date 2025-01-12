package job

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return
}
