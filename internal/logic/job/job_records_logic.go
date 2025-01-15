package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type JobRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobRecordsLogic {
	return &JobRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobRecordsLogic) JobRecords(req *types.JobRecordsRequest) (resp *types.JobRecordsResponse, err error) {
	// 查询job记录
	jobRecords, total, err := l.svcCtx.JobRecordModel.FindPage(l.ctx, req.Current, req.PageSize, req.JobId, req.StartTime, req.EndTime, req.Status)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "查询job记录失败")
	}

	lists := make([]types.JobRecords, len(jobRecords))
	for i, record := range jobRecords {
		lists[i] = types.JobRecords{
			JobId:    record.JobId,
			JobName:  record.JobName,
			ExecTime: utils.FormatDate(record.ExecTime),
			Status:   record.Status,
			TraceId:  record.TraceId,
			Param:    record.Param,
			Result:   record.Result,
		}
	}

	return &types.JobRecordsResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		List:     lists,
	}, nil
}
