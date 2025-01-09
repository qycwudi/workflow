package api

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/api"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ApisecretKeyUpdateExpirationTimeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApiSecretKeyUpdateExpirationTimeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := api.NewApisecretKeyUpdateExpirationTimeLogic(r.Context(), svcCtx)
		resp, err := l.ApisecretKeyUpdateExpirationTime(&req)
		response.Response(w, resp, err)

	}
}
