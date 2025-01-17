package model

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ModuleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModuleListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := model.NewModuleListLogic(r.Context(), svcCtx)
		resp, err := l.ModuleList(&req)
		response.Response(w, resp, err)

	}
}
