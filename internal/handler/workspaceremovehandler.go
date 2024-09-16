package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gogogo/internal/logic"
	"gogogo/internal/svc"
	"gogogo/internal/types"
)

func WorkSpaceRemoveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkRemoveRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewWorkSpaceRemoveLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceRemove(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
