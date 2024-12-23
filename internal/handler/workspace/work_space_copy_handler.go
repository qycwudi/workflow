package workspace

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"workflow/internal/logic/workspace"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func WorkSpaceCopyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkSpaceCopyRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := workspace.NewWorkSpaceCopyLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceCopy(&req)
		response.Response(w, resp, err)

	}
}
