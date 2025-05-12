package trace

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/pkg/engine"
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
	pool := engine.GetGlobalPool()
	resp = &types.QueryGoroutinesResponse{
		Cap:      pool.Cap(),
		Running:  pool.Running(),
		Free:     pool.Free(),
		Waiting:  pool.Waiting(),
		IsClosed: pool.IsClosed(),
	}
	return resp, nil
}
