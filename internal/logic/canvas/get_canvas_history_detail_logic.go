package canvas

import (
	"context"

	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetCanvasHistoryDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCanvasHistoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCanvasHistoryDetailLogic {
	return &GetCanvasHistoryDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCanvasHistoryDetailLogic) GetCanvasHistoryDetail(req *types.GetCanvasHistoryDetailReq) (resp *types.GetCanvasHistoryDetailResp, err error) {
	canvasHistory, err := l.svcCtx.CanvasHistoryModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布历史版本失败")
	}

	// 转map
	graph := make(map[string]interface{})
	err = json.Unmarshal([]byte(canvasHistory.Draft), &graph)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布历史版本失败")
	}

	return &types.GetCanvasHistoryDetailResp{
		Id:    canvasHistory.Id,
		Name:  canvasHistory.Name,
		Graph: graph,
	}, nil
}
