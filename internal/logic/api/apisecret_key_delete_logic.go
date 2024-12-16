package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApisecretKeyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApisecretKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApisecretKeyDeleteLogic {
	return &ApisecretKeyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApisecretKeyDeleteLogic) ApisecretKeyDelete(req *types.ApiSecretKeyDeleteRequest) (resp *types.ApiSecretKeyDeleteResponse, err error) {

	// 逻辑删除
	err = l.svcCtx.ApiSecretKeyModel.LogicalDelete(l.ctx, req.SecretKey)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除 API Secret Key 记录失败")
	}
	resp = &types.ApiSecretKeyDeleteResponse{
		SecretKey: req.SecretKey,
	}
	return resp, nil
}
