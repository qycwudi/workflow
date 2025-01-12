package canvas

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type SaveCanvasHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveCanvasHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveCanvasHistoryLogic {
	return &SaveCanvasHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveCanvasHistoryLogic) SaveCanvasHistory(req *types.SaveCanvasHistoryReq) (resp *types.SaveCanvasHistoryResp, err error) {
	// 查询当前画布草稿
	canvasDraft, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.WorkspaceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草稿失败")
	}
	if canvasDraft == nil {
		return nil, errors.New(int(logic.ParamError), "画布草稿不存在")
	}
	result, err := l.svcCtx.CanvasHistoryModel.Insert(l.ctx, &model.CanvasHistory{
		WorkspaceId: req.WorkspaceId,
		Draft:       canvasDraft.Draft,
		Name:        req.Name,
		CreateTime:  time.Now(),
		Mode:        model.CanvasHistoryModeDraft,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "保存画布历史版本失败")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "保存画布历史版本失败")
	}
	resp = &types.SaveCanvasHistoryResp{
		Id: id,
	}

	return resp, nil
}
