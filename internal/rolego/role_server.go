package rolego

import (
	"context"
	"github.com/rulego/rulego/api/types"
	endpoint2 "github.com/rulego/rulego/api/types/endpoint"
	"github.com/rulego/rulego/endpoint"
	"github.com/rulego/rulego/endpoint/rest"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"workflow/internal/utils"
)

type RoleServer struct {
}

// InitRoleServer 注册接口 注册规则链
func InitRoleServer() {
	config := types.Config{Logger: &utils.RoleCustomLog{}}
	restEndpoint, err := endpoint.Registry.New(rest.Type, config, rest.Config{Server: ":9999"})
	if err != nil {
		log.Fatal(err)
	}
	_, err = restEndpoint.AddRouter(Route(), "POST")
	if err != nil {
		logx.Errorf("Role Server Init AddRouter err: %v\n", err)
	}
	err = restEndpoint.Start()
	if err != nil {
		logx.Errorf("Role Server Init Start err: %v\n", err)
	}
	logx.Info("init role server success")

	logx.Info("init role server load api")
	apis, err := RoleChain.svc.ApiModel.FindByOn(context.Background())
	if err != nil {
		logx.Errorf("find api server error: %s\n", err.Error())
		return
	}
	for _, api := range apis {
		logx.Infof("load api id:%s,name:%s", api.ApiId, api.ApiName)
		RoleChain.LoadChain(api.ApiId, []byte(api.Dsl))
	}
	logx.Infof("init role server load api complete : %d", len(apis))
}

func Route() endpoint2.Router {
	// 如果需要把规则链执行结果同步响应给客户端，则增加wait语义
	router := endpoint.NewRouter().From("/api/role/v1/:chainId").
		To("chain:${chainId}").
		// 必须增加Wait，异步转同步，http才能正常响应，如果不响应同步响应，不要加这一句，会影响吞吐量
		Wait().
		Process(func(router endpoint2.Router, exchange *endpoint.Exchange) bool {
			err := exchange.Out.GetError()
			if err != nil {
				// 错误
				exchange.Out.SetStatusCode(400)
				exchange.Out.SetBody([]byte(exchange.Out.GetError().Error()))
			} else {
				// 把处理结果响应给客户端，http endpoint 必须增加 Wait()，否则无法正常响应
				outMsg := exchange.Out.GetMsg()
				exchange.Out.Headers().Set("Content-Type", "application/json")
				exchange.Out.SetBody([]byte(outMsg.Data))
			}
			return true
		}).End()
	return router
}
