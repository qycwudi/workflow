package canvas

import (
	"os"
	"testing"

	"workflow/internal/rulego"
)

func TestCanvasRunLogic_CanvasRun(t *testing.T) {
	// 读取文件
	file, err := os.ReadFile("/Users/qiangyuecheng/GolandProjects/work-flow/internal/logic/canvas/dsl/test-3.json")
	if err != nil {
		panic(err)
	}
	// 解析文件
	canvasId, ruleChain, err := rulego.ParsingDsl(string(file))
	if err != nil {
		t.Fatal("failed to parse dsl")
	}
	t.Logf("canvasId:%s", canvasId)
	t.Logf("ruleChain:%s", ruleChain)
	// 运行文件
	rulego.RoleChain.LoadChain(canvasId, ruleChain)
	matadata := make(map[string]string)
	matadata["name"] = "雪兔"
	matadata["age"] = "18"
	// 读取 start param
	data := "{\"name\":\"雪兔\",\"age\":18}"
	result := rulego.RoleChain.Run(canvasId, matadata, data)
	t.Logf("chain run result:%+v", result)
}

func TestCanvasRunLogic_RunAll(t *testing.T) {
	// 读取文件
	file, err := os.ReadFile("/Users/qiangyuecheng/GolandProjects/work-flow/internal/logic/canvas/dsl/test-5.json")
	if err != nil {
		panic(err)
	}
	// 解析文件
	canvasId, ruleChain, err := rulego.ParsingDsl(string(file))
	if err != nil {
		t.Fatal("failed to parse dsl")
	}
	t.Logf("canvasId:%s", canvasId)
	t.Logf("ruleChain:%s", ruleChain)
	// 运行文件
	rulego.RoleChain.LoadChain(canvasId, ruleChain)
	matadata := make(map[string]string)
	matadata["name"] = "雪兔"
	matadata["age"] = "18"
	// 读取 start param
	data := "{\"name\":\"雪兔\",\"age\":18}"
	result := rulego.RoleChain.Run(canvasId, matadata, data)
	t.Logf("chain run result:%+v", result)
}
