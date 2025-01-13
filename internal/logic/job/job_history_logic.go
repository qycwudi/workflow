package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type JobHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobHistoryLogic {
	return &JobHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobHistoryLogic) JobHistory(req *types.JobHistoryRequest) (resp *types.JobHistoryResponse, err error) {
	canvasHistory, total, err := l.svcCtx.CanvasHistoryModel.FindAllJobByWorkspaceId(l.ctx, req.WorkspaceId, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询历史版本失败")
	}

	lists := make([]types.JobHistory, len(canvasHistory))
	for i, history := range canvasHistory {
		lists[i] = types.JobHistory{
			Id:          history.Id,
			WorkspaceId: history.WorkspaceId,
			JobName:     history.Name,
			CreateTime:  history.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	resp = &types.JobHistoryResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		List:     lists,
	}

	return resp, nil
}
