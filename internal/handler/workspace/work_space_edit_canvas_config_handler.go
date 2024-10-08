package workspace

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/workspace"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func WorkSpaceEditCanvasConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkSpaceUploadCanvasConfigTagRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := workspace.NewWorkSpaceEditCanvasConfigLogic(r.Context(), svcCtx)
		resp, err := l.WorkSpaceEditCanvasConfig(&req)
		response.Response(w, resp, err)

	}
}
