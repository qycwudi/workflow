package kv

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/kv"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func CommonGetVByKHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetVByKRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := kv.NewCommonGetVByKLogic(r.Context(), svcCtx)
		resp, err := l.CommonGetVByK(&req)
		response.Response(w, resp, err)

	}
}
