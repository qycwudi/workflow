package rulego

import (
	"os"
	"testing"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/utils"
)

func Test_roleChain_Run_Join(t *testing.T) {
	file, _ := os.ReadFile("./chain/base.json")
	config := rulego.NewConfig()
	config.Logger = &utils.RoleCustomLog{}
	chain, err := rulego.New(
		"ctdbgr2flvkk9k4q7t7g",
		file,
		rulego.WithConfig(config),
		// types.WithAspects(&DebugAop{})
	)
	if err != nil {
		logx.Errorf("load role chain fail,err:%v\n", err)
		return
	}
	matadata := make(map[string]string)
	matadata["env"] = "jlhalsjdhfoisdbv"
	data := "{\"name\": \"张三\",  \"age\": 10}"

	msg := types.NewMsg(0, "CANVAS_MSG", types.JSON, matadata, data)
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
	}))
	t.Logf("chain run result:%+v", result)
}

func Test_roleChain_Run_Single(t *testing.T) {
	file, _ := os.ReadFile("./chain/join.json")
	config := rulego.NewConfig()
	config.Logger = &utils.RoleCustomLog{}
	chain, err := rulego.New(
		"cs8vfql3sjtkvhkubul",
		file,
		rulego.WithConfig(config),
		// types.WithAspects(&DebugAop{})
	)
	if err != nil {
		logx.Errorf("load role chain fail,err:%v\n", err)
		return
	}
	matadata := make(map[string]string)
	matadata["env"] = "jlhalsjdhfoisdbv"
	data := "[{\"name\": \"张三\",  \"age\": 10},{\"name\": \"李四\",  \"age\": 10}]"

	msg := types.NewMsg(0, "CANVAS_MSG", types.JSON, matadata, data)
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithTellNext("1729233809531"), types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
	}))
	t.Logf("chain run result:%+v", result)
}

func Test_roleChain_Run_For(t *testing.T) {
	file, _ := os.ReadFile("./chain/for.json")
	config := rulego.NewConfig()
	logConf := logx.LogConf{
		Encoding: "plain",
	}

	logx.SetUp(logConf)

	// config.Logger = &utils.RoleCustomLog{}
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
	matadata["env"] = "jlhalsjdhfoisdbv"
	data := "{\"slice\": [{\"name\": \"张三\",  \"age\": 10},{\"name\": \"李四\",  \"age\": 10}]}"

	msg := types.NewMsg(0, "CANVAS_MSG", types.JSON, matadata, data)
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
	}))
	t.Logf("chain run result:%+v", result)
	/*
		do = s3 表示 s3是处理迭代逻辑的节点
		执行完以后 s2 要连 s4继续往后走，要想达到我要的那种效果，把 s2的connections 删了就行了
		{
		        "fromId": "s2",
		        "toId": "s4",
		        "type": "Success"
		 }
	*/
}

func Test_roleChain_Run_SingleNode(t *testing.T) {
	file, _ := os.ReadFile("./chain/single.json")
	config := rulego.NewConfig()
	logConf := logx.LogConf{
		Encoding: "plain",
	}

	logx.SetUp(logConf)

	// config.Logger = &utils.RoleCustomLog{}
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
	matadata["env"] = "jlhalsjdhfoisdbv"
	data := "{\"slice\": [{\"name\": \"张三\",  \"age\": 10},{\"name\": \"李四\",  \"age\": 10}]}"

	msg := types.NewMsg(0, "CANVAS_MSG", types.JSON, matadata, data)
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
	}))
	t.Logf("chain run result:%+v", result)
	/*
		do = s3 表示 s3是处理迭代逻辑的节点
		执行完以后 s2 要连 s4继续往后走，要想达到我要的那种效果，把 s2的connections 删了就行了
		{
		        "fromId": "s2",
		        "toId": "s4",
		        "type": "Success"
		 }
	*/
}
