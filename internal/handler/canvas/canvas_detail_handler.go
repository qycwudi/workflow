package canvas

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/canvas"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func CanvasDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CanvasDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := canvas.NewCanvasDetailLogic(r.Context(), svcCtx)
		resp, err := l.CanvasDetail(&req)
		response.Response(w, resp, err)

	}
}
