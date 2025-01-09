package api

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/api"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ApiPublishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApiPublishRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := api.NewApiPublishLogic(r.Context(), svcCtx)
		resp, err := l.ApiPublish(&req)
		response.Response(w, resp, err)

	}
}
