package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiSecretyKeyUpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiSecretyKeyUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiSecretyKeyUpdateStatusLogic {
	return &ApiSecretyKeyUpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiSecretyKeyUpdateStatusLogic) ApiSecretyKeyUpdateStatus(req *types.ApiSecretyKeyUpdateStatusRequest) (resp *types.ApiSecretyKeyUpdateStatusResponse, err error) {

	// 检查下状态
	if req.Status != model.ApiSecretKeyStatusOn && req.Status != model.ApiSecretKeyStatusOff {
		return nil, errors.New(int(logic.SystemOrmError), "状态错误")
	}
	// 修改状态
	err = l.svcCtx.ApiSecretKeyModel.UpdateStatus(l.ctx, req.SecretKey, req.Status)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改 API Secret Key 状态失败")
	}
	resp = &types.ApiSecretyKeyUpdateStatusResponse{
		SecretKey: req.SecretKey,
		Status:    req.Status,
	}
	return resp, nil
}
