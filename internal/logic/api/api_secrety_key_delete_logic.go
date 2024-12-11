package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiSecretyKeyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiSecretyKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiSecretyKeyDeleteLogic {
	return &ApiSecretyKeyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiSecretyKeyDeleteLogic) ApiSecretyKeyDelete(req *types.ApiSecretyKeyDeleteRequest) (resp *types.ApiSecretyKeyDeleteResponse, err error) {

	// 逻辑删除
	err = l.svcCtx.ApiSecretKeyModel.LogicalDelete(l.ctx, req.SecretKey)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除 API Secret Key 记录失败")
	}
	resp = &types.ApiSecretyKeyDeleteResponse{
		SecretKey: req.SecretKey,
	}
	return resp, nil
}
