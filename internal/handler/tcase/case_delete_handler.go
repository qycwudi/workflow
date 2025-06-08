package tcase

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/tcase"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func CaseDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CaseDeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := tcase.NewCaseDeleteLogic(r.Context(), svcCtx)
		resp, err := l.CaseDelete(&req)
		response.Response(w, resp, err)

	}
}
