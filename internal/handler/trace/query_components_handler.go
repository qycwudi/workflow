package trace

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/trace"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func QueryComponentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryComponentsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := trace.NewQueryComponentsLogic(r.Context(), svcCtx)
		resp, err := l.QueryComponents(&req)
		response.Response(w, resp, err)

	}
}
