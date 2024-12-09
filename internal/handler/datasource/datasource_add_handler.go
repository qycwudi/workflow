package datasource

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/datasource"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func DatasourceAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DatasourceAddRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := datasource.NewDatasourceAddLogic(r.Context(), svcCtx)
		resp, err := l.DatasourceAdd(&req)
		response.Response(w, resp, err)

	}
}
