package rolego

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"testing"
	"workflow/internal/utils"
)

func Test_roleChain_Run(t *testing.T) {

	file, _ := os.ReadFile("/Users/qiangyuecheng/GolandProjects/work-flow/internal/rolego/chain/chain2.json")
	config := rulego.NewConfig()
	config.Logger = &utils.RoleCustomLog{}
	chain, err := rulego.New(
		"cs8vfql3sjtkvhkubul",
		file,
		rulego.WithConfig(config),
		types.WithAspects(&DebugAop{}))
	if err != nil {
		logx.Errorf("load role chain fail,err:%v\n", err)
		return
	}
	matadata := make(map[string]string)
	data := "{\"name\": \"张三\",  \"age\": 10}"

	msg := types.NewMsg(0, "CANVAS_MSG", types.JSON, matadata, data)
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
	}))
	t.Logf("chain run result:%+v", result)
}
