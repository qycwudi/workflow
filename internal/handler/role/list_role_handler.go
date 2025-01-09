package role

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/role"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func ListRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := role.NewListRoleLogic(r.Context(), svcCtx)
		resp, err := l.ListRole(&req)
		response.Response(w, resp, err)

	}
}
