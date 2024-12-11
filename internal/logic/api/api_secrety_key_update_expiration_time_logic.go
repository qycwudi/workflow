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

type ApiSecretyKeyUpdateExpirationTimeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiSecretyKeyUpdateExpirationTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiSecretyKeyUpdateExpirationTimeLogic {
	return &ApiSecretyKeyUpdateExpirationTimeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiSecretyKeyUpdateExpirationTimeLogic) ApiSecretyKeyUpdateExpirationTime(req *types.ApiSecretyKeyUpdateExpirationTimeRequest) (resp *types.ApiSecretyKeyUpdateExpirationTimeResponse, err error) {
	// 检查API是否存在
	_, err = l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "API 不存在")
	}

	// 检查过期时间
	if req.ExpirationTime <= time.Now().UnixMilli() {
		return nil, errors.New(int(logic.SystemOrmError), "过期时间不能小于当前时间")
	}

	// 修改过期时间
	expirationTime := time.UnixMilli(req.ExpirationTime)
	err = l.svcCtx.ApiSecretKeyModel.UpdateExpirationTime(l.ctx, req.ApiId, expirationTime)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改 API Secret Key 过期时间失败")
	}
	resp = &types.ApiSecretyKeyUpdateExpirationTimeResponse{
		ApiId:          req.ApiId,
		ExpirationTime: expirationTime.Format(time.DateTime),
	}
	return resp, nil
}
