package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"gogogo/internal/logic"
	"gogogo/internal/svc"
	"gogogo/internal/types"
	"gogogo/response"
	"net/http"
)

func WorkSpaceNewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkSpaceNewRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewWorkSpaceNewLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceNew(&req)
		response.Response(w, resp, err)

	}
}
