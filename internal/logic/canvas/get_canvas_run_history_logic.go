package canvas

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetCanvasRunHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCanvasRunHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCanvasRunHistoryLogic {
	return &GetCanvasRunHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCanvasRunHistoryLogic) GetCanvasRunHistory(req *types.GetCanvasRunHistoryReq) (resp *types.GetCanvasRunHistoryResp, err error) {
	l.Infof("GetCanvasRunHistory req:%+v", req)
	// 查询
	records, err := l.svcCtx.SpaceRecordModel.FindHistory(l.ctx, req.WorkSpaceId)
	if err != nil {
		l.Errorf("GetCanvasRunHistory err:%+v", err)
		return nil, errors.New(int(logic.SystemOrmError), "查询运行记录失败")
	}
	// 转换记录格式
	historyRecords := make([]types.RunHistoryRecord, 0, len(records))
	for _, record := range records {
		historyRecords = append(historyRecords, types.RunHistoryRecord{
			Id:             record.SerialNumber,
			StartTime:      record.RunTime.Format("2006-01-02 15:04:05"),
			Duration:       record.Duration,
			Status:         record.Status,
			ComponentCount: 0,
		})
	}

	resp = &types.GetCanvasRunHistoryResp{
		Records: historyRecords,
		Total:   int64(len(records)),
	}
	return resp, nil
}
