package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	asynq2 "github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"workflow/internal/asynq"
	"workflow/internal/cache"
	"workflow/internal/config"
	"workflow/internal/datasource"
	"workflow/internal/handler"
	"workflow/internal/locks"
	"workflow/internal/pubsub"
	"workflow/internal/rulego"
	"workflow/internal/svc"
)

var configFile = flag.String("f", "etc/workflow-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// 读取配置文件
	conf.MustLoad(*configFile, &c)

	// 从环境变量读取配置
	if port := os.Getenv("PORT"); port != "" {
		c.Port, _ = strconv.Atoi(port)
	}

	if apiPort := os.Getenv("API_PORT"); apiPort != "" {
		c.ApiPort, _ = strconv.Atoi(apiPort)
	}

	if mysqlDsn := os.Getenv("DSN"); mysqlDsn != "" {
		c.MySqlUrn = mysqlDsn
	}

	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		c.Redis.Host = redisHost
	}

	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		c.Redis.Password = redisPassword
	}

	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		c.Redis.DB, _ = strconv.Atoi(redisDB)
	}

	if logMode := os.Getenv("LOG_MODE"); logMode != "" {
		c.Log.Mode = logMode
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.Log.Level = logLevel
	}

	if ruleServerTrace := os.Getenv("RULE_SERVER_TRACE"); ruleServerTrace != "" {
		c.RuleServerTrace, _ = strconv.ParseBool(ruleServerTrace)
	}

	if ruleServerLimitSize := os.Getenv("RULE_SERVER_LIMIT_SIZE"); ruleServerLimitSize != "" {
		c.RuleServerLimitSize, _ = strconv.Atoi(ruleServerLimitSize)
	}

	// # 需要通过的域名，这里可以写多个域名 或者可以写 * 全部通过
	domains := []string{"*", "http://workflow", "http://127.0.0.1", "http://localhost"}
	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCors(domains...),
		rest.WithCustomCors(func(header http.Header) {
			// # 这里写允许通过的header key 不区分大小写
			header.Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id,OS,Platform, Version")
			header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
			header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		}, nil, "*"),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// 初始化 redis
	cache.NewRedis(ctx.RedisClient)
	// 注册规则链
	rulego.InitRoleChain(ctx)
	// 注册链服务
	rulego.InitRoleServer(c.RuleServerTrace, c.ApiPort, ctx.Config.RuleServerLimitSize)
	// 初始化锁
	locks.CustomLock = locks.NewLock("mysql", ctx)
	// 初始化数据源连接池
	datasource.InitDataSourceManager(ctx)
	// 初始化 asynq
	asynq.NewAsynqServer(ctx)
	// 初始化 asynq 周期性任务
	asynq.NewAsynqJob(ctx)
	// 初始化订阅
	pubsub.NewPubSub(ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in main:", r)
			}
		}()
		h := asynqmon.New(asynqmon.Options{
			RootPath: "/asynq",
			RedisConnOpt: asynq2.RedisClientOpt{
				Addr:     ctx.Config.Redis.Host,
				DB:       ctx.Config.Redis.DB,
				Username: "",
				Password: ctx.Config.Redis.Password,
			},
		})

		http.Handle(h.RootPath()+"/", h)
		println("Starting asynq monitor at 0.0.0.0:7201...")
		logx.Error(ctx, http.ListenAndServe(":7201", nil).Error())
	}()
	server.Start()

}
