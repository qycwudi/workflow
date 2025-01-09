package permission

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/permission"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func DeletePermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeletePermissionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := permission.NewDeletePermissionLogic(r.Context(), svcCtx)
		resp, err := l.DeletePermission(&req)
		response.Response(w, resp, err)

	}
}
