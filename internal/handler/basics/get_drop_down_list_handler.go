package basics

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/basics"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func GetDropDownListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDropDownListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := basics.NewGetDropDownListLogic(r.Context(), svcCtx)
		resp, err := l.GetDropDownList(&req)
		response.Response(w, resp, err)

	}
}
