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

type ApiSecretKeyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiSecretKeyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiSecretKeyListLogic {
	return &ApiSecretKeyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiSecretKeyListLogic) ApiSecretKeyList(req *types.ApiSecretKeyListRequest) (resp *types.ApiSecretKeyListResponse, err error) {
	total, secretKey, err := l.svcCtx.ApiSecretKeyModel.FindByApiIdPage(l.ctx, req.ApiId, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 API Secret Key 记录失败")
	}

	keys := make([]types.ApiSecretKey, len(secretKey))
	for i, key := range secretKey {
		keys[i] = types.ApiSecretKey{
			ApiId:          key.ApiId,
			Name:           key.Name,
			SecretKey:      key.SecretKey,
			ExpirationTime: utils.FormatDate(key.ExpirationTime),
			Status:         key.Status,
		}
	}

	resp = &types.ApiSecretKeyListResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		List:     keys,
	}
	return resp, nil
}
