package canvas

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/workflow"
)

type CanvasRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunLogic {
	return &CanvasRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunLogic) CanvasRun(req *types.CanvasRunRequest) (resp *types.CanvasRunResponse, err error) {

	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}

	err = workflow.Register(l.ctx, canvas.Draft)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "注册任务流失败")
	}

	// 读取 参数
	data := req.Params

	// 运行文件
	serialId, result, err := workflow.Run(l.ctx, canvas.WorkspaceId, data)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "运行任务流失败")
	}
	resp = &types.CanvasRunResponse{
		Id:     serialId,
		Data:   result.Output,
		Status: strings.Join(result.Route, ","),
		Error:  result.Error,
	}
	return resp, nil
}
