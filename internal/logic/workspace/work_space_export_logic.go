package workspace

import (
	"context"
	"encoding/base64"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type WorkSpaceExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceExportLogic {
	return &WorkSpaceExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceExportLogic) WorkSpaceExport(req *types.WorkSpaceExportRequest) (resp *types.WorkSpaceExportResponse, err error) {
	// 查询 canvas 信息
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "查询画布信息失败")
	}
	// 转成 base64
	export := base64.StdEncoding.EncodeToString([]byte(canvas.Draft))
	resp = &types.WorkSpaceExportResponse{
		Export: export,
	}
	return resp, nil
}
