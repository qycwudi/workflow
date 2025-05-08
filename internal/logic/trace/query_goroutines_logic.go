package trace

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryGoroutinesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryGoroutinesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryGoroutinesLogic {
	return &QueryGoroutinesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryGoroutinesLogic) QueryGoroutines(req *types.QueryGoroutinesRequest) (resp *types.QueryGoroutinesResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
