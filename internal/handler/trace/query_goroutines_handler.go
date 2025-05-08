package trace

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/trace"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func QueryGoroutinesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryGoroutinesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := trace.NewQueryGoroutinesLogic(r.Context(), svcCtx)
		resp, err := l.QueryGoroutines(&req)
		response.Response(w, resp, err)

	}
}
