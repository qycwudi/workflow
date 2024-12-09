package datasource

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DatasourceEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceEditLogic {
	return &DatasourceEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceEditLogic) DatasourceEdit(req *types.DatasourceEditRequest) (resp *types.DatasourceEditResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
