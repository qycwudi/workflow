package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type ApiSecretyKeyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiSecretyKeyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiSecretyKeyListLogic {
	return &ApiSecretyKeyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiSecretyKeyListLogic) ApiSecretyKeyList(req *types.ApiSecretyKeyListRequest) (resp *types.ApiSecretyKeyListResponse, err error) {
	total, secretKey, err := l.svcCtx.ApiSecretKeyModel.FindByApiIdPage(l.ctx, req.ApiId, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 API Secret Key 记录失败")
	}

	keys := make([]types.ApiSecretyKey, len(secretKey))
	for i, key := range secretKey {
		keys[i] = types.ApiSecretyKey{
			ApiId:          key.ApiId,
			SecretyKey:     key.SecretKey,
			ExpirationTime: utils.FormatDate(key.ExpirationTime),
		}
	}

	resp = &types.ApiSecretyKeyListResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		List:     keys,
	}
	return resp, nil
}
