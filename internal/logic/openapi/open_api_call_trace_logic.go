package openapi

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
)

type OpenApiCallTraceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type OpenApiCallTraceRequest struct {
	ApiId string         `path:"apiId"` // API ID
	Param map[string]any `json:"param"` // 参数
}

func NewOpenApiCallTraceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenApiCallTraceLogic {
	return &OpenApiCallTraceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenApiCallTraceLogic) OpenApiCallTrace(req *OpenApiCallTraceRequest) (resp map[string]any, err error) {
	// traceId := "trace-" + trace.TraceIDFromContext(l.ctx)
	return
}
