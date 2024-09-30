package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func WorkSpaceEditTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkSpaceEditTagRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewWorkSpaceEditTagLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceEditTag(&req)
		response.Response(w, resp, err)

	}
}
