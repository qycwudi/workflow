package canvas

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
	"workflow/internal/utils"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasRunRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunRecordLogic {
	return &CanvasRunRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunRecordLogic) CanvasRunRecord(req *types.CanvasRunRecordRequest) (resp *types.CanvasRunRecordResponse, err error) {
	canvasRunRecords, err := l.svcCtx.SpaceRecordModel.FindAll(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询运行记录失败")
	}
	records := make([]types.RunRecord, len(canvasRunRecords))
	for i, record := range canvasRunRecords {
		records[i] = types.RunRecord{
			TraceId: record.SerialNumber,
			Status:  record.Status,
			RunTime: utils.FormatDate(record.RunTime),
		}
	}
	resp = &types.CanvasRunRecordResponse{Records: records}
	return resp, nil
}
