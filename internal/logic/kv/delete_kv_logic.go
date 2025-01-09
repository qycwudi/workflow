package kv

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type DeleteKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKvLogic {
	return &DeleteKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteKvLogic) DeleteKv(req *types.DeleteKvRequest) (resp *types.DeleteKvResponse, err error) {
	// 查询key是否存在
	kv, err := l.svcCtx.KvModel.FindOneByKey(l.ctx, req.Key)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	if kv == nil {
		return nil, errors.New("key not found")
	}

	err = l.svcCtx.KvModel.Delete(l.ctx, kv.Id)
	if err != nil {
		return nil, err
	}
	return &types.DeleteKvResponse{}, nil
}
