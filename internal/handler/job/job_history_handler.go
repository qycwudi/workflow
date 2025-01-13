package job

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"workflow/internal/logic/job"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func JobHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JobHistoryRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := job.NewJobHistoryLogic(r.Context(), svcCtx)
		resp, err := l.JobHistory(&req)
		response.Response(w, resp, err)

	}
}
