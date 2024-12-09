package datasource

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
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
	// 查询数据源是否存在
	datasource, err := l.svcCtx.DatasourceModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "数据源不存在")
	}

	// 更新数据源信息
	datasource.Type = req.Type
	datasource.Config = req.Config
	datasource.Switch = int64(req.Switch)
	datasource.Hash = req.Hash
	datasource.Status = req.Status
	datasource.UpdateTime = time.Now()

	err = l.svcCtx.DatasourceModel.Update(l.ctx, datasource)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "修改数据源失败")
	}

	resp = &types.DatasourceEditResponse{
		Id: req.Id,
	}
	return resp, nil
}
