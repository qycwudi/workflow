package canvas

import (
	"context"
	"encoding/json"
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
	canvasDraft, err := json.Marshal(req.CanvasDraft)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "保存画布历史版本失败")
	}
	result, err := l.svcCtx.CanvasHistoryModel.Insert(l.ctx, &model.CanvasHistory{
		WorkspaceId: req.WorkspaceId,
		Draft:       string(canvasDraft),
		Name:        req.Name,
		CreateTime:  time.Now(),
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
