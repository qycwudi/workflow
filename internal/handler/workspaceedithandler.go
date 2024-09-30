package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func WorkSpaceEditHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkSpaceEditRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewWorkSpaceEditLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceEdit(&req)
		response.Response(w, resp, err)
	}
}
