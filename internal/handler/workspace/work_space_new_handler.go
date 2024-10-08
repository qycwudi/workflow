package workspace

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/workspace"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func WorkSpaceNewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkSpaceNewRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := workspace.NewWorkSpaceNewLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceNew(&req)
		response.Response(w, resp, err)

	}
}
