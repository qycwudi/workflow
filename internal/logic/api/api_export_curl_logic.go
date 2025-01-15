package api

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiExportCurlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiExportCurlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiExportCurlLogic {
	return &ApiExportCurlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiExportCurlLogic) ApiExportCurl(req *types.ApiExportCurlRequest) (resp *types.ApiExportCurlResponse, err error) {
	// 构建基础的 curl 命令
	curl := "curl -X POST"

	// 添加 URL
	curl += " '" + req.Url + "'"

	// 处理请求头
	var headers map[string]string
	if err = json.Unmarshal([]byte(req.Header), &headers); err != nil {
		return nil, errors.New(int(logic.SystemError), "解析请求头失败")
	}

	// 检查是否存在 Content-Type 头
	hasContentType := false
	for k, v := range headers {
		if strings.ToLower(k) == "content-type" {
			hasContentType = true
		}
		curl += " -H '" + k + ": " + v + "'"
	}

	// 如果没有 Content-Type 头,默认添加 application/json
	if !hasContentType {
		curl += " -H 'Content-Type: application/json'"
	}

	// 添加请求体
	if req.Body != "" {
		curl += " -d '" + req.Body + "'"
	}

	resp = &types.ApiExportCurlResponse{
		Curl: curl,
	}
	return resp, nil
}
