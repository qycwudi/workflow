package datasource

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DatasourceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceListLogic {
	return &DatasourceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceListLogic) DatasourceList(req *types.DatasourceListRequest) (resp *types.DatasourceListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
