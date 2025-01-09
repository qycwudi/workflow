package canvas

import (
	"context"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	util "workflow/internal/utils"
)

type RestoreCanvasHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreCanvasHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreCanvasHistoryLogic {
	return &RestoreCanvasHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestoreCanvasHistoryLogic) RestoreCanvasHistory(req *types.RestoreCanvasHistoryReq) (resp *types.RestoreCanvasHistoryResp, err error) {
	canvasHistory, err := l.svcCtx.CanvasHistoryModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布历史版本失败")
	}
	// 查询当前画布草稿
	canvasDraft, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, canvasHistory.WorkspaceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草稿失败")
	}
	// 更新画布
	canvasDraft.Draft = canvasHistory.Draft
	canvasDraft.UpdateAt = time.Now()
	userId, err := util.GetUserId(l.ctx)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "获取用户id失败")
	}
	userIdStr := strconv.FormatInt(userId, 10)
	canvasDraft.UpdateBy = userIdStr
	err = l.svcCtx.CanvasModel.Update(l.ctx, canvasDraft)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "恢复画布历史版本失败")
	}
	resp = &types.RestoreCanvasHistoryResp{
		Id:          canvasHistory.Id,
		WorkspaceId: canvasHistory.WorkspaceId,
	}
	return resp, nil
}
