package api

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApisecretKeyUpdateExpirationTimeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApisecretKeyUpdateExpirationTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApisecretKeyUpdateExpirationTimeLogic {
	return &ApisecretKeyUpdateExpirationTimeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApisecretKeyUpdateExpirationTimeLogic) ApisecretKeyUpdateExpirationTime(req *types.ApiSecretKeyUpdateExpirationTimeRequest) (resp *types.ApiSecretKeyUpdateExpirationTimeResponse, err error) {

	// 检查过期时间
	if req.ExpirationTime <= time.Now().UnixMilli() {
		return nil, errors.New(int(logic.SystemOrmError), "过期时间不能小于当前时间")
	}

	// 修改过期时间
	expirationTime := time.UnixMilli(req.ExpirationTime)
	err = l.svcCtx.ApiSecretKeyModel.UpdateExpirationTime(l.ctx, req.SecretKey, expirationTime)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改 API Secret Key 过期时间失败")
	}
	resp = &types.ApiSecretKeyUpdateExpirationTimeResponse{
		ExpirationTime: expirationTime.Format(time.DateTime),
	}
	return resp, nil
}
