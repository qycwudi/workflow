package canvas

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/tidwall/gjson"
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
	// 读取请求参数
	draft := gjson.Parse(canvas.Draft)
	startNode := draft.Get("graph.nodes").Array()[0]
	code := startNode.Get("data.code")
	if code.String() == "" {
		return nil, errors.New(int(logic.SystemError), "开始节点运行参数为空")
	}
	param := gjson.Parse(code.String())
	metaData := param.Get("metaData").String()
	if metaData == "" {
		return nil, errors.New(int(logic.SystemError), "开始节点运行参数 metaData 为空")
	}
	metadata := make(map[string]string)
	_ = json.Unmarshal([]byte(metaData), &metadata)
	data := param.Get("data").String()
	if data == "" {
		return nil, errors.New(int(logic.SystemError), "开始节点运行参数 data 为空")
	}
	// 运行文件
	rolego.RoleChain.LoadChain(canvasId, ruleChain)
	result := rolego.RoleChain.Run(canvasId, metadata, data)
	l.Infof("chain run result:%+v", result)

	respData := make(map[string]interface{})

	err = json.Unmarshal([]byte(result.Data), &respData)
	resp = &types.CanvasRunResponse{
		Id:       result.Id,
		Ts:       result.Ts,
		MetaData: result.Metadata,
		Data:     respData,
	}
	return
}
