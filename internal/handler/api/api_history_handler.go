package api

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/api"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ApiHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApiHistoryRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := api.NewApiHistoryLogic(r.Context(), svcCtx)
		resp, err := l.ApiHistory(&req)
		response.Response(w, resp, err)

	}
}