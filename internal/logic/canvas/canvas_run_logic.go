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

	// 解析文件
	// 1. 读取基本信息
	canvasJson := gjson.ParseBytes([]byte(canvas.Draft))
	canvasId := canvasJson.Get("id").String()
	l.Infof("canvasId:%s\n", canvasId)

	// 2. 构造点
	graphNodes := canvasJson.Get("graph.nodes").Array()
	nodes := make([]rolego.Node, len(graphNodes))
	for i, node := range graphNodes {
		r := rolego.Node{
			Id:   node.Get("id").String(),
			Type: node.Get("data.type").String(),
			Name: node.Get("data.title").String(),
			// 不同组件配置读取逻辑不同
			Configuration: rolego.ModuleReadConfig(node.Get("data")),
		}
		nodes[i] = r
	}
	for _, node := range nodes {
		l.Infof("%+v\n", node)
	}

	// 3. 构造线
	graphEdges := canvasJson.Get("graph.edges").Array()
	edges := make([]rolego.Connection, len(graphEdges))
	for i, edge := range graphEdges {
		r := rolego.Connection{
			FromId: edge.Get("source").String(),
			ToId:   edge.Get("target").String(),
			Type:   edge.Get("relation").String(),
		}
		edges[i] = r
	}
	for _, edge := range edges {
		l.Infof("%+v\n", edge)
	}

	// 4. 构造执行实体
	ruleChain := rolego.Rule{
		RuleChain: rolego.RuleChain{Id: canvasId},
		Metadata: rolego.Metadata{
			Nodes:       nodes,
			Connections: edges,
		},
	}
	//
	ruleChainMar, err := json.Marshal(ruleChain)
	if err != nil {
		panic(err)
	}
	l.Infof("%s\n", string(ruleChainMar))
	// 运行文件
	rolego.RoleChain.LoadChain(canvasId, ruleChainMar)
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
