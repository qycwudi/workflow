package user

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/user"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func UserUpdateInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUpdateInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewUserUpdateInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserUpdateInfo(&req)
		response.Response(w, resp, err)

	}
}
