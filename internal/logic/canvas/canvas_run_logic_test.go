package canvas

import (
	"github.com/rulego/rulego/utils/json"
	"github.com/tidwall/gjson"
	"os"
	"testing"
	"workflow/internal/rolego"
)

func TestCanvasRunLogic_CanvasRun(t *testing.T) {
	// 读取文件
	file, err := os.ReadFile("/Users/qiangyuecheng/GolandProjects/work-flow/internal/logic/canvas/dsl/test-1.json")
	if err != nil {
		panic(err)
	}
	// 解析文件
	// 1. 读取基本信息
	canvasJson := gjson.ParseBytes(file)
	canvasId := canvasJson.Get("id").String()
	t.Logf("canvasId:%s\n", canvasId)

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
	t.Logf("node: \n")
	for _, node := range nodes {
		t.Logf("%+v\n", node)
	}
	t.Logf("\n")

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
	t.Logf("edge: \n")
	for _, edge := range edges {
		t.Logf("%+v\n", edge)
	}
	t.Logf("\n")

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
	t.Logf("%s\n", string(ruleChainMar))
	// 运行文件
	rolego.RoleChain.LoadChain(canvasId, ruleChainMar)
	matadata := make(map[string]string)
	data := "{\"name\":\"雪兔\",\"age\":18}"
	result := rolego.RoleChain.Run(canvasId, matadata, data)
	t.Logf("chain run result:%+v", result)
}
