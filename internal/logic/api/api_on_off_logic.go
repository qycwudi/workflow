package api

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiOnOffLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiOnOffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiOnOffLogic {
	return &ApiOnOffLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiOnOffLogic) ApiOnOff(req *types.ApiOnOffRequest) (resp *types.ApiOnOffResponse, err error) {
	err = l.svcCtx.ApiModel.UpdateStatusByApiId(l.ctx, req.ApiId, req.Status)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改API发布状态失败")
	}
	return &types.ApiOnOffResponse{
		ApiId:  req.ApiId,
		Status: req.Status,
	}, nil
}
