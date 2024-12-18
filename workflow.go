package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"workflow/internal/config"
	"workflow/internal/corn"
	"workflow/internal/datasource"
	"workflow/internal/handler"
	"workflow/internal/locks"
	"workflow/internal/rulego"
	"workflow/internal/svc"
)

var configFile = flag.String("f", "etc/workflow-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// 读取配置文件
	conf.MustLoad(*configFile, &c)

	// 从环境变量读取 服务端口
	if port := os.Getenv("PORT"); port != "" {
		c.Port, _ = strconv.Atoi(port)
	}

	if apiPort := os.Getenv("API_PORT"); apiPort != "" {
		c.ApiPort, _ = strconv.Atoi(apiPort)
	}

	// 从环境变量读取 mysql 配置
	if mysqlDsn := os.Getenv("DSN"); mysqlDsn != "" {
		c.MySqlUrn = mysqlDsn
	}

	// 从环境变量读取 log 配置
	if logMode := os.Getenv("LOG_MODE"); logMode != "" {
		c.Log.Mode = logMode
	}

	// 从环境变量读取 log 配置
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.Log.Level = logLevel
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
	// 注册规则链
	rulego.InitRoleChain(ctx)
	// 注册链服务
	rulego.InitRoleServer(c.ApiPort)
	// 初始化锁
	locks.CustomLock = locks.NewLock("mysql", ctx)
	// 初始化Job
	corn.NewJob(c.Job, ctx)
	// 初始化数据源连接池
	datasource.InitDataSourceManager(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
