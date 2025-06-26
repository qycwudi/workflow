package basics

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/basics"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func GetHomeStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetHomeStatisticsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := basics.NewGetHomeStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.GetHomeStatistics(&req)
		response.Response(w, resp, err)

	}
}
