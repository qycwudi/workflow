package kv

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSetKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommonSetKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSetKvLogic {
	return &CommonSetKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommonSetKvLogic) CommonSetKv(req *types.SetKvRequest) (resp *types.SetKvResponse, err error) {

	response := &types.SetKvResponse{}

	// 查询key是否存在
	value, err := l.svcCtx.GoGoGoKvModel.FindByKey(l.ctx, req.Key)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = l.svcCtx.GoGoGoKvModel.Insert(l.ctx, &model.GogogoKv{
			SpiderName: req.SpiderName,
			K:          req.Key,
			V:          req.Value,
			Timestamp:  time.Now().UnixMilli(),
		})
		if err != nil {
			response.Code = int(logic.SystemError)
			response.Message = fmt.Sprintf("insert key error %v", err.Error())
			return response, nil
		}
		response.Code = int(logic.SUCCESS)
		response.Message = "SUCCESS"
		return response, nil
	}

	if err != nil {
		response.Code = int(logic.SystemError)
		response.Message = fmt.Sprintf("find key error %v", err)
	}
	if value != nil {
		response.Code = int(logic.KeyExist)
		response.Message = fmt.Sprintf("key exist :%+v", value)
	}
	return response, nil
}