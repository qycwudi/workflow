package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"workflow/internal/config"
	"workflow/internal/corn"
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
	rulego.InitRoleServer()
	// 初始化锁
	locks.CustomLock = locks.NewLock("mysql", ctx)
	// 初始化Job
	corn.NewJob(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
