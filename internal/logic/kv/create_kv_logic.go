package kv

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type CreateKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateKvLogic {
	return &CreateKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateKvLogic) CreateKv(req *types.CreateKvRequest) (resp *types.CreateKvResponse, err error) {
	// 查询key是否存在
	kv, err := l.svcCtx.KvModel.FindOneByKey(l.ctx, req.Key)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	if kv != nil {
		return nil, errors.New("key already exists")
	}

	kv = &model.Kv{
		Key:   req.Key,
		Value: req.Value,
	}

	_, err = l.svcCtx.KvModel.Insert(l.ctx, kv)
	if err != nil {
		return nil, err
	}
	return &types.CreateKvResponse{}, nil
}
