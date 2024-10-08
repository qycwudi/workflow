package kv

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"workflow/internal/logic"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonGetVByKLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommonGetVByKLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonGetVByKLogic {
	return &CommonGetVByKLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommonGetVByKLogic) CommonGetVByK(req *types.GetVByKRequest) (resp *types.GetVByKResponse, err error) {
	response := types.GetVByKResponse{}

	value, err := l.svcCtx.GoGoGoKvModel.FindByKey(l.ctx, req.Key)
	if errors.Is(err, sqlx.ErrNotFound) {
		response.Code = int(logic.KeyMiss)
		response.Message = "miss key " + req.Key
		return &response, nil
	}

	if err != nil {
		response.Code = int(logic.SystemError)
		response.Message = "system error" + err.Error()
		return &response, nil
	}

	response.Code = int(logic.SUCCESS)
	response.Message = "hit key"
	response.Value = value.V
	return &response, nil
}
