package rolego

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
	"workflow/internal/utils"
)

//
// func InitRoleChain() {
// 	file, _ := os.ReadFile("/Users/qiangyuecheng/GolandProjects/work-flow/internal/rolego/chain/chain1.json")
// 	config := rulego.NewConfig()
// 	config.Logger = &utils.RoleCustomLog{}
// 	_, err := rulego.New(
// 		"rule8848",
// 		file,
// 		rulego.WithConfig(config),
// 		types.WithAspects(&Trace{}))
// 	if err != nil {
// 		logx.Errorf("init role chain fail,err:%v\n", err)
// 		return
// 	}
// 	logx.Info("init role chain success")
// 	//
// 	for i := 0; i < 5; i++ {
// 		go func(index int) {
// 			get, _ := rulego.Get("rule01")
// 			js := "{\"userName\":\"xue\",\"password\":\"123456\",\"role\":\"admin2\"}"
// 			metaData := types.NewMetadata()
// 			metaData.PutValue("userName", "string")
// 			metaData.PutValue("password", "string")
// 			metaData.PutValue("role", "string"+strconv.Itoa(index))
// 			msg := types.NewMsg(0, "TELEMETRY_MSG", types.JSON, metaData, js)
// 			get.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
// 				println("GetSelfId:", ctx.GetSelfId())
// 				println("msgId:", msg.Id)
// 			}))
// 		}(i)
// 	}
//
// 	time.Sleep(5 * time.Second)
// }

var RoleChain = new(roleChain)

type roleChain struct {
}

func (r *roleChain) LoadChain(id string, json []byte) {
	chain, b := rulego.Get(id)
	if b {
		// 重新加载
		err := chain.ReloadSelf(json)
		if err != nil {
			logx.Errorf("reload self role chain %s fail,err:%v\n", id, err)
		}
		logx.Infof("reload self role chain %s success,json: %s \n", id, string(json))
		return
	}
	config := rulego.NewConfig()
	config.Logger = &utils.RoleCustomLog{}
	_, err := rulego.New(
		id,
		json,
		rulego.WithConfig(config),
		types.WithAspects(&Trace{}))
	if err != nil {
		logx.Errorf("load role chain fail,err:%v\n", err)
		return
	}
	logx.Infof("load %s role chain success,json:%s \n", id, string(json))
}

func (r *roleChain) getChain(id string) types.RuleEngine {
	chain, b := rulego.Get(id)
	if !b {
		logx.Errorf("get role chain fail,id:%s\n", id)
	}
	return chain
}

func (r *roleChain) Run(id string, metadata map[string]string, data string) types.RuleEngine {
	logx.Infof("id:%s\n metadata:%+v\n data:%s\n", id, metadata, data)
	chain, b := rulego.Get(id)
	if !b {
		logx.Errorf("get role chain fail,id:%s\n", id)
	}
	metaData := types.NewMetadata()
	for k, v := range metadata {
		metaData.PutValue(k, v)
	}
	msg := types.NewMsg(0, "TELEMETRY_MSG", types.JSON, metaData, data)
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		println("Run-GetSelfId:", ctx.GetSelfId())
		println("Run-msgId:", msg.Id)
	}))
	return chain
}
