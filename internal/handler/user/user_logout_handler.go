package user

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/user"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func UserLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLogoutRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewUserLogoutLogic(r.Context(), svcCtx)
		resp, err := l.UserLogout(&req)
		response.Response(w, resp, err)

	}
}
