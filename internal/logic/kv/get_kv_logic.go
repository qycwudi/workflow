package kv

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetKvLogic {
	return &GetKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKvLogic) GetKv(req *types.GetKvRequest) (resp *types.GetKvResponse, err error) {
	// 查询key是否存在
	kv, err := l.svcCtx.KvModel.FindOneByKey(l.ctx, req.Key)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	if kv == nil {
		return nil, errors.New("key not found")
	}

	return &types.GetKvResponse{
		Kv: types.Kv{
			Key:   kv.Key,
			Value: kv.Value,
		},
	}, nil
}
