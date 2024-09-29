package rolego

import (
	"github.com/rulego/rulego"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func InitRoleChain() {
	file, _ := os.ReadFile("/Users/qiangyuecheng/GolandProjects/work-flow/internal/rolego/chain/chain1.json")
	_, err := rulego.New("rule8848", file)
	if err != nil {
		logx.Errorf("init role chain fail,err:%v\n", err)
		return
	}
	logx.Info("init role chain success")
	//
	// for i := 0; i < 5; i++ {
	// 	go func(index int) {
	// 		get, _ := rulego.Get("rule01")
	// 		js := "{\"userName\":\"xue\",\"password\":\"123456\",\"role\":\"admin2\"}"
	// 		metaData := types.NewMetadata()
	// 		metaData.PutValue("userName", "string")
	// 		metaData.PutValue("password", "string")
	// 		metaData.PutValue("role", "string"+strconv.Itoa(index))
	// 		msg := types.NewMsg(0, "TELEMETRY_MSG", types.JSON, metaData, js)
	// 		get.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
	// 			println("GetSelfId:", ctx.GetSelfId())
	// 			println("msgId:", msg.Id)
	// 		}))
	// 	}(i)
	// }
	//
	// time.Sleep(5 * time.Second)
}
