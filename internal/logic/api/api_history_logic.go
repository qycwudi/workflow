package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiHistoryLogic {
	return &ApiHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiHistoryLogic) ApiHistory(req *types.ApiHistoryRequest) (resp *types.ApiHistoryResponse, err error) {
	canvasHistory, total, err := l.svcCtx.CanvasHistoryModel.FindAllApiByWorkspaceId(l.ctx, req.WorkspaceId, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询历史版本失败")
	}

	lists := make([]types.ApiHistory, len(canvasHistory))
	for i, history := range canvasHistory {
		lists[i] = types.ApiHistory{
			Id:          history.Id,
			WorkspaceId: history.WorkspaceId,
			Name:        history.Name,
			CreateTime:  history.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	resp = &types.ApiHistoryResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		List:     lists,
	}

	return resp, nil
}
