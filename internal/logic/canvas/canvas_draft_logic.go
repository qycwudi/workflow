package canvas

import (
	"context"
	errors2 "errors"
	"strconv"
	"time"

	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/utils"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	util "workflow/internal/utils"
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
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	draftMarshal, _ := json.Marshal(req)
	userId, err := util.GetUserId(l.ctx)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "获取用户id失败")
	}
	userIdStr := strconv.FormatInt(userId, 10)
	if errors2.Is(err, sqlc.ErrNotFound) {
		// 新增
		_, err = l.svcCtx.CanvasModel.Insert(l.ctx, &model.Canvas{
			WorkspaceId: req.Id,
			Draft:       string(draftMarshal),
			CreateAt:    time.Now(),
			UpdateAt:    time.Now(),
			CreateBy:    userIdStr,
			UpdateBy:    userIdStr,
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemOrmError), "新增画布草案失败")
		}
		return nil, nil
	}
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	// 更新
	canvas.Draft = string(draftMarshal)
	canvas.UpdateAt = time.Now()
	canvas.UpdateBy = userIdStr
	err = l.svcCtx.CanvasModel.Update(l.ctx, canvas)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "更新画布草案失败")
	}
	resp = &types.CanvasDraftResponse{
		Hash:       utils.NewUuid(),
		UpdateTime: time.Now().UnixMilli(),
	}

	return resp, nil
}
