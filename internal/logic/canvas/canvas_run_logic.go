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
	// 1. 读取点、线
	// 查询点
	nodes, err := l.svcCtx.NodeModel.FindOneByWorkSpace(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询节点失败")
	}

	// 查询边
	edges, err := l.svcCtx.EdgeModel.FindOneByWorkSpace(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询边失败")
	}

	// 2. 拼接 json
	rule := rolego.Rule{
		RuleChain: rolego.RuleChain{Id: req.Id},
		Metadata: rolego.Metadata{
			Nodes:       make([]rolego.Node, 0, len(nodes)),
			Connections: make([]rolego.Connection, 0, len(edges)),
		},
	}

	// 将nodes转换为Rule中的Node格式
	for _, dbNode := range nodes {
		node := rolego.Node{
			Id:            dbNode.NodeId,
			Type:          dbNode.NodeType,
			Name:          dbNode.NodeName,
			Configuration: make(map[string]interface{}),
		}
		// 假设这里需要将Configuration字段转换为map[string]interface{}
		if err := json.Unmarshal([]byte(dbNode.Configuration), &node.Configuration); err != nil {
			return nil, errors.New(int(logic.SystemOrmError), "节点配置解析失败")
		}
		rule.Metadata.Nodes = append(rule.Metadata.Nodes, node)
	}

	// 将edges转换为Rule中的Connection格式
	for _, dbEdge := range edges {
		connection := rolego.Connection{
			FromId: dbEdge.Source,
			ToId:   dbEdge.Target,
			Type:   dbEdge.Route,
		}
		rule.Metadata.Connections = append(rule.Metadata.Connections, connection)
	}

	ruleJson, err := json.Marshal(rule)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "规则序列化失败")
	}
	l.Info(string(ruleJson))
	// 3. 加载到链池 记录 md5 新建 or 重新加载
	rolego.RoleChain.LoadChain(req.Id, ruleJson)
	// 4. doMsg
	dataMar, _ := json.Marshal(req.Data)
	rolego.RoleChain.Run(req.Id, req.MetaData, string(dataMar))

	resp = &types.CanvasRunResponse{
		Id:       "",
		Response: string(ruleJson),
	}
	return
}
