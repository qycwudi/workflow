package api

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"workflow/internal/logic/api"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ApiListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApiPublishListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := api.NewApiListLogic(r.Context(), svcCtx)
		resp, err := l.ApiList(&req)
		response.Response(w, resp, err)
	}
}
