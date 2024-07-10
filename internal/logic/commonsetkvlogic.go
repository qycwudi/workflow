package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gogogo/internal/model"
	"gogogo/internal/svc"
	"gogogo/internal/types"
	"time"

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
	value, err := l.svcCtx.GogogoKvModel.FindByKey(l.ctx, req.Key)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = l.svcCtx.GogogoKvModel.Insert(l.ctx, &model.GogogoKv{
			SpiderName: req.SpiderName,
			K:          req.Key,
			V:          req.Value,
			Timestamp:  int64(time.Now().Second()),
		})
		if err != nil {
			response.Code = int(SystemError)
			response.Message = fmt.Sprintf("insert key error %v", err.Error())
			return response, nil
		}
		response.Code = int(SUCCESS)
		response.Message = "SUCCESS"
		return response, nil
	}

	if err != nil {
		response.Code = int(SystemError)
		response.Message = fmt.Sprintf("find key error %v", err)
	}
	if value != nil {
		response.Code = int(KeyExist)
		response.Message = fmt.Sprintf("key exist :%+v", value)
	}
	return response, nil
}
