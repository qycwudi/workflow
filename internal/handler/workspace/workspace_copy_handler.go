package workspace

import (
	"net/http"
	"workflow/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"workflow/internal/logic/workspace"
	"workflow/internal/svc"
	"workflow/internal/types"
)

// WorkspaceCopyHandler 画布复制
func WorkspaceCopyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkspaceCopyRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := workspace.NewWorkspaceCopyLogic(r.Context(), svcCtx)
		resp, err := l.WorkspaceCopy(&req)
		response.Response(w, resp, err)
	}
}
