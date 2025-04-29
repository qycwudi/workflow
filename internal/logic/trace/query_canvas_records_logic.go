package trace

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type QueryCanvasRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询画布运行记录列表
func NewQueryCanvasRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCanvasRecordsLogic {
	return &QueryCanvasRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryCanvasRecordsLogic) QueryCanvasRecords(req *types.QueryCanvasRecordsRequest) (resp *types.QueryCanvasRecordsResponse, err error) {
	// 查询画布运行记录
	total, records, err := l.svcCtx.SpaceRecordModel.FindByWorkspaceIdPage(l.ctx, req.Id, req.Current, req.PageSize)
	if err != nil {
		l.Errorf("query canvas records err:%+v", err)
		return nil, errors.New(int(logic.SystemOrmError), "查询画布运行记录失败")
	}

	// 转换画布运行记录
	canvasRecords := make([]types.CanvasRecord, len(records))
	for i, record := range records {
		canvasRecords[i] = types.CanvasRecord{
			Id:        record.WorkspaceId,
			SerialId:  record.SerialNumber,
			Status:    record.Status,
			Duration:  record.Duration,
			StartTime: utils.FormatDate(record.RunTime),
			Data:      record.Other,
		}
	}
	resp = &types.QueryCanvasRecordsResponse{
		Total:    total,
		Current:  req.Current,
		PageSize: req.PageSize,
		Records:  canvasRecords,
	}

	return resp, nil
}
