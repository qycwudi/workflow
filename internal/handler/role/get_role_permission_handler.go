package role

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/role"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func GetRolePermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRolePermissionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := role.NewGetRolePermissionLogic(r.Context(), svcCtx)
		resp, err := l.GetRolePermission(&req)
		response.Response(w, resp, err)

	}
}
