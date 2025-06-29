package canvas

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/utils"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	util "workflow/internal/utils"
	pkgutils "workflow/pkg/utils"
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
	draftMarshal, _ := sonic.Marshal(req.Graph)
	userId, _ := util.GetUserId(l.ctx)
	userIdStr := strconv.FormatInt(userId, 10)
	
	// 记录画布草稿操作的审计日志
	pkgutils.LogBusinessOperation(pkgutils.OpUpdate, req.Id, userIdStr, 
		logx.Field("operation", "canvas_draft_save"),
		logx.Field("draft_size", len(draftMarshal)))
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
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
				pkgutils.LogModuleError(pkgutils.ModuleCanvas, "canvas_draft_creation", err, logx.Field("workspace_id", req.Id))
				return nil, errors.New(int(logic.SystemOrmError), "Failed to add canvas draft")
			}
			resp = &types.CanvasDraftResponse{
				Hash:       utils.NewUuid(),
				UpdateTime: time.Now().UnixMilli(),
			}
			return resp, nil
		} else {
			pkgutils.LogModuleError(pkgutils.ModuleCanvas, "canvas_draft_query", err, logx.Field("workspace_id", req.Id))
			return nil, errors.New(int(logic.SystemOrmError), "Failed to query canvas draft")
		}
	}
	// 更新
	canvas.Draft = string(draftMarshal)
	canvas.UpdateAt = time.Now()
	canvas.UpdateBy = userIdStr
	err = l.svcCtx.CanvasModel.Update(l.ctx, canvas)
	if err != nil {
		pkgutils.LogModuleError(pkgutils.ModuleCanvas, "canvas_draft_update", err, logx.Field("workspace_id", req.Id))
		return nil, errors.New(int(logic.SystemOrmError), "Failed to update canvas draft")
	}
	resp = &types.CanvasDraftResponse{
		Hash:       utils.NewUuid(),
		UpdateTime: time.Now().UnixMilli(),
	}
	return resp, nil
}
