package datasource

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
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
	err = l.svcCtx.DatasourceModel.Delete(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "删除数据源失败")
	}
	resp = &types.DatasourceDeleteResponse{
		Id: req.Id,
	}
	return resp, nil
}
