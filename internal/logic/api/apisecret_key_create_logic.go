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

type ApisecretKeyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApisecretKeyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApisecretKeyCreateLogic {
	return &ApisecretKeyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApisecretKeyCreateLogic) ApisecretKeyCreate(req *types.ApiSecretKeyCreateRequest) (resp *types.ApiSecretKeyCreateResponse, err error) {
	// 检查API是否存在
	_, err = l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "API 不存在")
	}

	// 检查过期时间
	if req.ExpirationTime <= time.Now().UnixMilli() {
		return nil, errors.New(int(logic.ParamError), "过期时间不能小于当前时间")
	}

	var secretKey string
	// 检查密钥是否存在
	if req.SecretKey != "" {
		// 校验长度
		if len(req.SecretKey) != 32 {
			return nil, errors.New(int(logic.ParamError), "密钥长度错误,必须为32位")
		}
		_, err = l.svcCtx.ApiSecretKeyModel.FindOneByApiIdAndSecretKey(l.ctx, req.ApiId, req.SecretKey)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, errors.New(int(logic.SystemOrmError), "密钥已存在")
			} else {
				return nil, errors.New(int(logic.SystemOrmError), "密钥查询失败")
			}
		}
	} else {
		// 生成密钥
		secretKey = utils.GenerateUUID()
	}

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

	resp = &types.ApiSecretKeyCreateResponse{
		ApiId:          req.ApiId,
		Name:           req.Name,
		SecretKey:      secretKey,
		ExpirationTime: expirationTime.Format(time.DateTime),
	}
	return resp, nil
}
