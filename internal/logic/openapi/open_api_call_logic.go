package openapi

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"

	"workflow/internal/svc"
	"workflow/internal/workflow"
)

type OpenApiCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type OpenApiCallRequest struct {
	ApiId string         `path:"apiId"` // API ID
	Param map[string]any `json:"param"` // 参数
}

func NewOpenApiCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenApiCallLogic {
	return &OpenApiCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenApiCallLogic) OpenApiCall(req *OpenApiCallRequest) (resp map[string]any, err error) {
	traceId := trace.TraceIDFromContext(l.ctx)
	workspaceId := req.ApiId
	param := req.Param
	// 执行
	_, result, err := workflow.Run(l.ctx, traceId, workspaceId, param)
	if err != nil {
		logx.Errorw("[画布] 执行工作流失败",
			logx.Field("工作空间ID", workspaceId),
			logx.Field("序列ID", traceId),
			logx.Field("错误", err))
		return
	}
	// 返回结果
	return result.Output.(map[string]any), nil
}
