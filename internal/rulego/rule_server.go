package rulego

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/rulego/rulego/api/types"
	endpoint2 "github.com/rulego/rulego/api/types/endpoint"
	"github.com/rulego/rulego/endpoint"
	"github.com/rulego/rulego/endpoint/rest"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/cache"
	"workflow/internal/model"
	"workflow/internal/utils"
)

type RoleServer struct {
}

// InitRoleServer 注册接口 注册规则链
func InitRoleServer(trace bool, apiPort int, limitSize int) {
	config := types.Config{Logger: &utils.RoleCustomLog{}}
	restEndpoint, err := endpoint.Registry.New(rest.Type, config, rest.Config{Server: fmt.Sprintf(":%d", apiPort)})
	if err != nil {
		log.Fatal(err)
	}
	restEndpoint.AddInterceptors(func(router endpoint2.Router, exchange *endpoint.Exchange) bool {
		startTime := time.Now()
		defer func() {
			endTime := time.Now()
			duration := endTime.Sub(startTime)
			logx.Infof("请求耗时: %v毫秒", duration.Milliseconds())
		}()

		// 检查请求体大小  1*1024*1024  1MB
		if len(exchange.In.Body()) > limitSize*1024*1024 { // 1MB = 1024 * 1024 bytes
			response := make(map[string]interface{})
			response["code"] = 413
			response["message"] = "状态请求实体太大"
			marshal, _ := json.Marshal(response)
			exchange.Out.SetStatusCode(http.StatusRequestEntityTooLarge)
			exchange.Out.SetBody(marshal)
			return false
		}

		// 检查 api 状态和 token
		chainId := exchange.In.GetParam("chainId")
		token := exchange.In.Headers().Get("Authorization")
		if token == "" {
			exchange.Out.SetStatusCode(http.StatusUnauthorized)
			exchange.Out.SetBody([]byte("token is required"))
			return false
		}
		token = strings.TrimPrefix(token, "Bearer ")

		checkStartTime := time.Now()
		if err := checkApiAndToken(chainId, token); err != nil {
			exchange.Out.SetStatusCode(http.StatusUnauthorized)
			exchange.Out.SetBody([]byte(err.Error()))
			return false
		}
		checkDuration := time.Since(checkStartTime)
		logx.Infof("检查API和Token耗时: %v毫秒", checkDuration.Milliseconds())
		// 设置metadata
		env := getApiEnvCache(chainId)
		for k, v := range env {
			exchange.In.GetMsg().Metadata[k] = v
		}
		exchange.In.GetMsg().Metadata["secret_key"] = token
		exchange.In.GetMsg().Metadata["api_id"] = chainId
		exchange.In.GetMsg().Metadata["api_name"] = chainId
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
		RoleChain.LoadApiServiceChain(api.ApiId, []byte(api.Dsl))
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
				if exchange.Out.GetMsg().Metadata.GetValue("END_NODE") == "true" {
					outMsg := exchange.Out.GetMsg()
					exchange.Out.Headers().Set("Content-Type", "application/json")
					exchange.Out.SetBody([]byte(outMsg.Data))
				}
			}
			return true
		}).End()
	return router
}

// 读取环境变量缓存
func getApiEnvCache(chainId string) map[string]string {
	var env map[string]string
	value, err := cache.Redis.Get(context.Background(), fmt.Sprintf(cache.EnvRedisKey, chainId))
	if err != nil {
		if err == red.Nil {
			logx.Errorf("Create new API environment variable cache %v", err)
			api, err := RoleChain.svc.ApiModel.FindOneByApiId(context.Background(), chainId)
			if err != nil {
				logx.Errorf("Failed to query API from database: %v", err)
			}
			// 查询API的环境变量
			workspace, err := RoleChain.svc.WorkSpaceModel.FindOneByWorkspaceId(context.Background(), api.WorkspaceId)
			if err != nil {
				logx.Errorf("Failed to query workspace from database: %v", err)
			}
			err = json.Unmarshal([]byte(workspace.Configuration), &env)
			if err != nil {
				logx.Errorf("Failed to parse workspace configuration: %v", err)
			}
			// 缓存环境变量
			err = cache.Redis.Set(context.Background(), fmt.Sprintf(cache.EnvRedisKey, chainId), workspace.Configuration, time.Hour*24)
			if err != nil {
				logx.Errorf("Failed to set API environment variable cache: %v", err)
			}
		} else {
			logx.Errorf("Failed to read API environment variable from Redis cache: %v", err)
		}
	} else {
		err = json.Unmarshal([]byte(value), &env)
		if err != nil {
			logx.Errorf("Failed to parse API environment variable from Redis cache: %v", err)
		}
	}
	return env
}

// 检查api和token
func checkApiAndToken(apiId, secretKey string) error {
	redisKey := fmt.Sprintf(cache.ApiSecretKeyRedisKey, apiId, secretKey)
	value, err := cache.Redis.Get(context.Background(), redisKey)

	// 如果Redis查询出错(非key不存在)
	if err != nil && err != red.Nil {
		logx.Errorf("get api and token error: %s", err.Error())
		return errors.New("get api and token error")
	}

	// key不存在,从数据库查询并缓存
	if err == red.Nil {
		return checkAndCacheApiAndToken(apiId, secretKey, redisKey)
	}

	// Redis中存在,检查状态
	return validateApiAndTokenStatus(value)
}

func checkAndCacheApiAndToken(apiId, secretKey, redisKey string) error {
	logx.Info("api and token not exists in cache, query from db")

	// 检查API
	api, err := RoleChain.svc.ApiModel.FindOneByApiId(context.Background(), apiId)
	if err != nil {
		logx.Errorf("api not found: %s", err.Error())
		_ = cache.Redis.Set(context.Background(), redisKey, "api_not_found", time.Hour*24)
		return errors.New("api not found")
	}
	if api.Status != model.ApiStatusOn {
		logx.Errorf("api status is off")
		_ = cache.Redis.Set(context.Background(), redisKey, "api_off", time.Hour*24)
		return errors.New("api status is off")
	}

	// 检查Token
	token, err := RoleChain.svc.ApiSecretKeyModel.FindOneByApiIdAndSecretKey(context.Background(), apiId, secretKey)
	if err != nil {
		logx.Errorf("token not found: %s", err.Error())
		_ = cache.Redis.Set(context.Background(), redisKey, "token_not_found", time.Hour*24)
		return errors.New("token not found")
	}
	if token.Status != model.ApiSecretKeyStatusOn {
		logx.Errorf("token status is off")
		_ = cache.Redis.Set(context.Background(), redisKey, "token_off", time.Hour*24)
		return errors.New("token status is off")
	}

	// 缓存过期时间
	value := fmt.Sprintf("valid:%d", token.ExpirationTime.Unix())
	_ = cache.Redis.Set(context.Background(), redisKey, value, time.Hour*24)
	return nil
}

func validateApiAndTokenStatus(value string) error {
	switch value {
	case "api_not_found":
		return errors.New("api not found")
	case "api_off":
		return errors.New("api status is off")
	case "token_not_found":
		return errors.New("token not found")
	case "token_off":
		return errors.New("token status is off")
	}

	// 检查token过期时间
	if strings.HasPrefix(value, "valid:") {
		expireStr := strings.TrimPrefix(value, "valid:")
		expirationTime, _ := strconv.ParseInt(expireStr, 10, 64)
		if time.Now().Unix() > expirationTime {
			return errors.New("token is expired")
		}
		return nil
	}

	return errors.New("invalid cache value")
}
