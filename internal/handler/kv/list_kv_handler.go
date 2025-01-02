package kv

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/kv"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ListKvHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListKvRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := kv.NewListKvLogic(r.Context(), svcCtx)
		resp, err := l.ListKv(&req)
		response.Response(w, resp, err)

	}
}
