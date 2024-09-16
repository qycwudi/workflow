package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gogogo/internal/logic"
	"gogogo/internal/svc"
	"gogogo/internal/types"
)

func WorkSpaceNewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkSpaceRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewWorkSpaceNewLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceNew(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
