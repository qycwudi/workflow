package canvas

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/canvas"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func CanvasEditNodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CanvasEditNodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := canvas.NewCanvasEditNodeLogic(r.Context(), svcCtx)
		resp, err := l.CanvasEditNode(&req)
		response.Response(w, resp, err)

	}
}
