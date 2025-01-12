package job

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/job"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func JobPublishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JobPublishRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := job.NewJobPublishLogic(r.Context(), svcCtx)
		resp, err := l.JobPublish(&req)
		response.Response(w, resp, err)

	}
}
