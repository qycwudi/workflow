package kv

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UpdateKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateKvLogic {
	return &UpdateKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateKvLogic) UpdateKv(req *types.UpdateKvRequest) (resp *types.UpdateKvResponse, err error) {
	// 查询key是否存在
	kv, err := l.svcCtx.KvModel.FindOneByKey(l.ctx, req.Key)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	if kv == nil {
		return nil, errors.New("key not found")
	}

	kv.Value = req.Value
	err = l.svcCtx.KvModel.Update(l.ctx, kv)
	if err != nil {
		return nil, err
	}

	return &types.UpdateKvResponse{}, nil
}
