package datasource

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/datasource"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func DatasourceDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DatasourceDeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := datasource.NewDatasourceDeleteLogic(r.Context(), svcCtx)
		resp, err := l.DatasourceDelete(&req)
		response.Response(w, resp, err)

	}
}
