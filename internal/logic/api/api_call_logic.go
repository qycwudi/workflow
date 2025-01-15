package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCallLogic {
	return &ApiCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCallLogic) ApiCall(req *types.ApiCallRequest) (resp *types.ApiCallResponse, err error) {
	// 检查请求头
	if req.Header == "" {
		return nil, errors.New(int(logic.SystemError), "请求头不能为空")
	}
	// 检查请求体
	if req.Body == "" {
		return nil, errors.New(int(logic.SystemError), "请求体不能为空")
	}
	// 检查请求体是否为json
	var jsonData interface{}
	if err := json.Unmarshal([]byte(req.Body), &jsonData); err != nil {
		return nil, errors.New(int(logic.SystemError), "请求体必须是有效的JSON格式")
	}
	// 检查url
	if req.Url == "" {
		return nil, errors.New(int(logic.SystemError), "url不能为空")
	}

	// 根据 ApiId 查询 API 信息
	_, err = l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 API 信息失败")
	}

	// 发送 HTTP 请求
	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", req.Url, strings.NewReader(req.Body))
	if err != nil {
		return nil, err
	}

	// 设置请求头
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
		httpReq.Header.Set(k, v)
	}

	// 如果没有 Content-Type 头,默认添加 application/json
	if !hasContentType {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// 执行请求
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "执行请求失败")
	}
	defer httpResp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "读取响应失败")
	}

	resp = &types.ApiCallResponse{
		Response: string(respBody),
	}
	return resp, nil
}
