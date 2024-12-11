package api

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type ApiSecretyKeyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiSecretyKeyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiSecretyKeyCreateLogic {
	return &ApiSecretyKeyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiSecretyKeyCreateLogic) ApiSecretyKeyCreate(req *types.ApiSecretyKeyCreateRequest) (resp *types.ApiSecretyKeyCreateResponse, err error) {
	// 检查API是否存在
	_, err = l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "API 不存在")
	}

	// 检查过期时间
	if req.ExpirationTime <= time.Now().UnixMilli() {
		return nil, errors.New(int(logic.SystemOrmError), "过期时间不能小于当前时间")
	}

	// 生成密钥
	secretKey := utils.GenerateUUID()
	expirationTime := time.UnixMilli(req.ExpirationTime)
	var apiSecretKey = model.ApiSecretKey{
		ApiId:          req.ApiId,
		Name:           req.Name,
		SecretKey:      secretKey,
		Status:         model.ApiSecretKeyStatusOn,
		IsDeleted:      0,
		ExpirationTime: expirationTime,
	}

	_, err = l.svcCtx.ApiSecretKeyModel.Insert(l.ctx, &apiSecretKey)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "生成 API Secret Key 失败")
	}

	resp = &types.ApiSecretyKeyCreateResponse{
		ApiId:          req.ApiId,
		Name:           req.Name,
		SecretyKey:     secretKey,
		ExpirationTime: expirationTime.Format(time.DateTime),
	}
	return resp, nil
}
