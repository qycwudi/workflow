package job

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"workflow/internal/logic/job"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/response"
)

func JobRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JobRecordsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := job.NewJobRecordsLogic(r.Context(), svcCtx)
		resp, err := l.JobRecords(&req)
		response.Response(w, resp, err)

	}
}
