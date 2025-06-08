package tcase

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/tcase"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func CaseListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CaseListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := tcase.NewCaseListLogic(r.Context(), svcCtx)
		resp, err := l.CaseList(&req)
		response.Response(w, resp, err)

	}
}
