package kv

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type ListKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListKvLogic {
	return &ListKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListKvLogic) ListKv(req *types.ListKvRequest) (resp *types.ListKvResponse, err error) {
	kvList, total, err := l.svcCtx.KvModel.FindAll(l.ctx, req.Key, req.Current, req.PageSize)
	if err != nil {
		return nil, err
	}

	kvs := make([]types.Kv, len(kvList))
	for i, kv := range kvList {
		kvs[i] = types.Kv{
			Key:   kv.Key,
			Value: kv.Value,
		}
	}

	return &types.ListKvResponse{
		Total: total,
		List:  kvs,
	}, nil
}
