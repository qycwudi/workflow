package rulego

import (
	"github.com/rulego/rulego/utils/json"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
)

func ParsingDsl(draft string) (string, []byte, error) {
	// 解析文件
	// 1. 读取基本信息
	canvasJson := gjson.ParseBytes([]byte(draft))
	canvasId := canvasJson.Get("id").String()
	logx.Infof("canvasId:%s\n", canvasId)

	// 特殊关系[迭代的处理节点]
	baseInfo := make(map[string]string, 0)

	// 2. 构造线
	graphEdges := canvasJson.Get("graph.edges").Array()
	edges := make([]Connection, len(graphEdges))
	for i, edge := range graphEdges {
		if edge.Get("sourceHandle").String() == "Do" {
			baseInfo[edge.Get("source").String()] = edge.Get("target").String()
			continue
		}
		r := Connection{
			FromId: edge.Get("source").String(),
			ToId:   edge.Get("target").String(),
			Type:   edge.Get("sourceHandle").String(),
		}
		edges[i] = r
	}

	// 3. 构造点
	graphNodes := canvasJson.Get("graph.nodes").Array()
	nodes := make([]Node, len(graphNodes))
	for i, node := range graphNodes {
		// 存储id等组件原始信息
		baseInfo["id"] = node.Get("id").String()
		r := Node{
			Id:   node.Get("id").String(),
			Type: node.Get("data.type").String(),
			Name: node.Get("data.name").String(),
			// 不同组件配置读取逻辑不同
			Configuration: ModuleReadConfig(node.Get("data"), baseInfo),
		}
		nodes[i] = r
	}

	// 4. 读取开始节点索引
	firstNodeIndex := 0
	for i, node := range nodes {
		if node.Type == Start {
			firstNodeIndex = i
			break
		}
	}

	// 4. 构造执行实体
	ruleChain := Rule{
		RuleChain: RuleChain{Id: canvasId},
		Metadata: Metadata{
			FirstNodeIndex: firstNodeIndex,
			Nodes:          nodes,
			Connections:    edges,
		},
	}

	ruleChainMar, err := json.Marshal(ruleChain)
	if err != nil {
		return "", nil, err
	}
	logx.Infof("%s", string(ruleChainMar))
	return canvasId, ruleChainMar, nil
}
