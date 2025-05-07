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
		logx.Errorw("[画布] 获取历史版本失败",
			logx.Field("历史版本ID", req.Id),
			logx.Field("错误", err))
		return nil, errors.New(int(logic.SystemOrmError), "获取历史版本失败")
	}
	// 查询当前画布草稿
	canvasDraft, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, canvasHistory.WorkspaceId)
	if err != nil {
		logx.Errorw("[画布] 获取工作流定义失败",
			logx.Field("工作空间ID", canvasHistory.WorkspaceId),
			logx.Field("错误", err))
		return nil, errors.New(int(logic.SystemOrmError), "获取工作流定义失败")
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
		logx.Errorw("[画布] 更新工作流定义失败",
			logx.Field("工作空间ID", canvasHistory.WorkspaceId),
			logx.Field("错误", err))
		return nil, errors.New(int(logic.SystemOrmError), "更新工作流定义失败")
	}
	resp = &types.RestoreCanvasHistoryResp{
		Id:          canvasHistory.Id,
		WorkspaceId: canvasHistory.WorkspaceId,
	}
	return resp, nil
}
