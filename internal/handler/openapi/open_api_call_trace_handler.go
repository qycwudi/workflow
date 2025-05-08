package openapi

import (
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/rest/httpx"

	"workflow/internal/logic/openapi"
	"workflow/internal/svc"
	"workflow/response"
)

func OpenApiCallTraceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req openapi.OpenApiCallTraceRequest
		// 解析 API ID
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		// 读取参数为 map[string]any
		var param map[string]any
		// 自定义实现读取 body 到 param
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		err = sonic.Unmarshal(body, &param)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		req.Param = param
		l := openapi.NewOpenApiCallTraceLogic(r.Context(), svcCtx)
		resp, err := l.OpenApiCallTrace(&req)
		response.Response(w, resp, err)
	}
}
