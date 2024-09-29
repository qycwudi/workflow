package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"gogogo/internal/config"
	"gogogo/internal/handler"
	"gogogo/internal/rolego"
	"gogogo/internal/svc"
	"gogogo/internal/utils"
	"time"
)

var configFile = flag.String("f", "etc/gogogo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	httpWriter := utils.NewHTTPLogWriter(c.OpenOB.Path, c.OpenOB.UserName, c.OpenOB.Password, 5*time.Second)
	logx.AddWriter(logx.NewWriter(httpWriter))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// 注册规则链
	rolego.InitRoleServer()
	rolego.InitRoleChain()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
