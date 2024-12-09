package datasource

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DatasourceDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceDeleteLogic {
	return &DatasourceDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceDeleteLogic) DatasourceDelete(req *types.DatasourceDeleteRequest) (resp *types.DatasourceDeleteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
