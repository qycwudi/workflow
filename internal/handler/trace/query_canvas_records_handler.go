package trace

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/trace"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func QueryCanvasRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryCanvasRecordsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := trace.NewQueryCanvasRecordsLogic(r.Context(), svcCtx)
		resp, err := l.QueryCanvasRecords(&req)
		response.Response(w, resp, err)

	}
}
