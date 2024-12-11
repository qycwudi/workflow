package api

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/api"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ApiSecretyKeyUpdateExpirationTimeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApiSecretyKeyUpdateExpirationTimeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := api.NewApiSecretyKeyUpdateExpirationTimeLogic(r.Context(), svcCtx)
		resp, err := l.ApiSecretyKeyUpdateExpirationTime(&req)
		response.Response(w, resp, err)

	}
}
