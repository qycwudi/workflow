package model

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ModuleNewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModuleNewRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := model.NewModuleNewLogic(r.Context(), svcCtx)
		resp, err := l.ModuleNew(&req)
		response.Response(w, resp, err)

	}
}
