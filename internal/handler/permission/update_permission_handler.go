package permission

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/permission"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func UpdatePermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePermissionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := permission.NewUpdatePermissionLogic(r.Context(), svcCtx)
		resp, err := l.UpdatePermission(&req)
		response.Response(w, resp, err)

	}
}
