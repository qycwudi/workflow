package api

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/api"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ApiSecretyKeyUpdateStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApiSecretyKeyUpdateStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := api.NewApiSecretyKeyUpdateStatusLogic(r.Context(), svcCtx)
		resp, err := l.ApiSecretyKeyUpdateStatus(&req)
		response.Response(w, resp, err)

	}
}
