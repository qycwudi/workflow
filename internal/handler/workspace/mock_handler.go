package workspace

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"workflow/internal/logic/workspace"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func MockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MockRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		// 打印请求头
		logx.Infof("request headers: %+v", r.Header)
		l := workspace.NewMockLogic(r.Context(), svcCtx)
		resp, err := l.Mock(&req)
		response.Response(w, resp, err)

	}
}
