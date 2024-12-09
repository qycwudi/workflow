package datasource

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/datasource"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func DatasourceListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DatasourceListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := datasource.NewDatasourceListLogic(r.Context(), svcCtx)
		resp, err := l.DatasourceList(&req)
		response.Response(w, resp, err)

	}
}
