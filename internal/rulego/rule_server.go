package rulego

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/rulego/rulego/api/types"
	endpoint2 "github.com/rulego/rulego/api/types/endpoint"
	"github.com/rulego/rulego/endpoint"
	"github.com/rulego/rulego/endpoint/rest"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
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
	restEndpoint.AddInterceptors(func(router endpoint2.Router, exchange *endpoint.Exchange) bool {

		// 检查请求体大小  1*1024*1024  1MB
		if len(exchange.In.Body()) > 1*1024*1024 { // 1MB = 1024 * 1024 bytes
			response := make(map[string]interface{})
			response["code"] = 413
			response["message"] = "状态请求实体太大"
			marshal, _ := json.Marshal(response)
			exchange.Out.SetStatusCode(http.StatusRequestEntityTooLarge)
			exchange.Out.SetBody(marshal)
			return false
		}

		// 检查 token
		token := exchange.In.Headers().Get("Authorization")
		if token == "" {
			exchange.Out.SetStatusCode(http.StatusUnauthorized)
			exchange.Out.SetBody([]byte("token is required"))
			return false
		}

		// 检查 api 状态
		chainId := exchange.In.GetParam("chainId")
		var apiName string
		if api, err := checkApi(chainId); err != nil {
			exchange.Out.SetStatusCode(http.StatusNotFound)
			exchange.Out.SetBody([]byte(err.Error()))
			return false
		} else {
			apiName = api.ApiName
		}

		// 检查 token 是否正确
		token = strings.TrimPrefix(token, "Bearer ")
		if err := checkToken(chainId, token); err != nil {
			exchange.Out.SetStatusCode(http.StatusUnauthorized)
			exchange.Out.SetBody([]byte(err.Error()))
			return false
		}
		exchange.In.GetMsg().Metadata["secret_key"] = token
		exchange.In.GetMsg().Metadata["api_id"] = chainId
		exchange.In.GetMsg().Metadata["api_name"] = apiName
		return true
	})
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

// 检查token
func checkToken(chainId, secretKey string) error {
	api, err := RoleChain.svc.ApiSecretKeyModel.FindOneBySecretKey(context.Background(), chainId, secretKey)
	if err != nil {
		return errors.New("this token is not found")
	}
	if api.Status != model.ApiSecretKeyStatusOn {
		return errors.New("this token status is off")
	}
	return nil
}

// 检查api是否存在且状态正常
func checkApi(apiId string) (model.Api, error) {
	api, err := RoleChain.svc.ApiModel.FindOneByApiId(context.Background(), apiId)
	if err != nil {
		return model.Api{}, errors.New("api not found")
	}
	if api.Status != model.ApiStatusOn {
		return model.Api{}, errors.New("api status is off")
	}
	return *api, nil
}
