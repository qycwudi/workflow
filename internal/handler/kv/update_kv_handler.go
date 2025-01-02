package kv

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/kv"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func UpdateKvHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateKvRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := kv.NewUpdateKvLogic(r.Context(), svcCtx)
		resp, err := l.UpdateKv(&req)
		response.Response(w, resp, err)

	}
}
