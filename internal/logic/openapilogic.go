package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenApiLogic {
	return &OpenApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenApiLogic) OpenApi(req *types.OpenApiRequest) (resp *types.OpenApiResponse, err error) {
	l.Infof("req:", req)
	response := map[string]string{}
	response["status"] = "SUCCESS"
	resp = &types.OpenApiResponse{
		Data: response,
	}
	return resp, nil
}
