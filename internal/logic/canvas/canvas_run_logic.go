package canvas

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
	"workflow/internal/rolego"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

	canvasId, ruleChain, err := rolego.ParsingDsl(canvas.Draft)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "解析画布草案失败")
	}
	// 运行文件
	rolego.RoleChain.LoadChain(canvasId, []byte(ruleChain))
	dataMar, _ := json.Marshal(req.Data)
	result := rolego.RoleChain.Run(canvasId, req.MetaData, string(dataMar))
	l.Infof("chain run result:%+v", result)

	respData := make(map[string]interface{})

	_ = json.Unmarshal([]byte(result.Data), &respData)

	resp = &types.CanvasRunResponse{
		Id:       result.Id,
		Ts:       result.Ts,
		MetaData: result.Metadata,
		Data:     respData,
	}
	return
}
