package rulego

import (
	"fmt"
	"github.com/rulego/rulego/utils/json"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strconv"
)

func ParsingDsl(draft string) (string, []byte, error) {
	// 处理全局变量
	draft = replaceGlobalVariables(draft)
	canvasJson := gjson.ParseBytes([]byte(draft))

	// 解析文件
	// 1. 读取基本信息
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
	firstNodeIndex := 0
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

		// 读取开始节点索引
		if r.Type == Start {
			firstNodeIndex = i
		}

		nodes[i] = r
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

func replaceGlobalVariables(draft string) string {
	globalMap := map[string]interface{}{}
	for _, global := range gjson.Get(draft, "globalParams").Array() {
		key := global.Get("key").String()
		value := global.Get("value").Value()
		globalMap[key] = value
	}

	re := regexp.MustCompile(`\${globalParams\.([\w_]+)}`)
	return re.ReplaceAllStringFunc(draft, func(match string) string {
		// 提取变量名
		varName := re.FindStringSubmatch(match)[1]

		// 从映射中获取变量值
		// 从映射中获取变量值
		if value, exists := globalMap[varName]; exists {
			// 根据类型进行相应的转换
			switch v := value.(type) {
			case string:
				return v
			case int:
				return strconv.Itoa(v)
			case float64:
				return strconv.FormatFloat(v, 'f', -1, 64)
			case bool:
				return strconv.FormatBool(v)
			default:
				// 对于其他类型，可以使用 fmt.Sprint
				return fmt.Sprint(v)
			}
		}

		// 如果变量不存在，可以选择返回原始字符串或空字符串
		return match
	})
}
