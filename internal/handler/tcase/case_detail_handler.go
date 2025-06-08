package tcase

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/tcase"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func CaseDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CaseDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := tcase.NewCaseDetailLogic(r.Context(), svcCtx)
		resp, err := l.CaseDetail(&req)
		response.Response(w, resp, err)

	}
}
