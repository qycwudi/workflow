package canvas

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetCanvasHistoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCanvasHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCanvasHistoryListLogic {
	return &GetCanvasHistoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCanvasHistoryListLogic) GetCanvasHistoryList(req *types.GetCanvasHistoryListReq) (resp *types.GetCanvasHistoryListResp, err error) {
	canvasHistoryList, err := l.svcCtx.CanvasHistoryModel.FindAll(l.ctx, &model.CanvasHistory{
		WorkspaceId: req.WorkspaceId,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取画布历史版本列表失败")
	}

	records := make([]types.CanvasHistoryRecord, 0)
	for _, v := range canvasHistoryList {
		records = append(records, types.CanvasHistoryRecord{
			Id:         v.Id,
			CreateTime: v.CreateTime.Format("2006-01-02 15:04:05"),
			Name:       v.Name,
		})
	}
	resp = &types.GetCanvasHistoryListResp{
		Records: records,
		Total:   int64(len(canvasHistoryList)),
	}
	return
}
