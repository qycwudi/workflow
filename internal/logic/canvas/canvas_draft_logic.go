package canvas

import (
	"context"
	errors2 "errors"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"
	"time"
	"workflow/internal/logic"
	"workflow/internal/model"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasDraftLogic {
	return &CanvasDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasDraftLogic) CanvasDraft(req *types.CanvasDraftRequest) (resp *types.CanvasDraftResponse, err error) {
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	draftMarshal, _ := json.Marshal(req)
	if errors2.Is(err, sqlc.ErrNotFound) {
		// 新增
		l.svcCtx.CanvasModel.Insert(l.ctx, &model.Canvas{
			WorkspaceId: req.Id,
			Draft:       string(draftMarshal),
			CreateAt:    time.Now(),
			UpdateAt:    time.Now(),
			CreateBy:    "admin",
			UpdateBy:    "admin",
		})
		return nil, nil
	}
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	// 更新
	canvas.Draft = string(draftMarshal)
	canvas.UpdateAt = time.Now()
	err = l.svcCtx.CanvasModel.Update(l.ctx, canvas)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "更新画布草案失败")
	}
	return nil, nil
}
