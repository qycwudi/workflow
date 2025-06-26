package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiGetApiDocLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiGetApiDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiGetApiDocLogic {
	return &ApiGetApiDocLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiGetApiDocLogic) ApiGetApiDoc(req *types.ApiGetApiDocRequest) (resp *types.ApiGetApiDocResponse, err error) {
	apiEntity, err := l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "api not found")
	}
	if apiEntity == nil {
		return nil, errors.New(int(logic.SystemOrmError), "api not found")
	}
	resp = &types.ApiGetApiDocResponse{
		ApiDoc: apiEntity.ApiDoc,
	}
	return resp, nil
}
