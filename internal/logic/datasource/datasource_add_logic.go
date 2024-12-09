package datasource

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DatasourceAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceAddLogic {
	return &DatasourceAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceAddLogic) DatasourceAdd(req *types.DatasourceAddRequest) (resp *types.DatasourceAddResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
