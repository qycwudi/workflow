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
	"os"
	"strconv"
	"time"
)

var configFile = flag.String("f", "etc/gogogo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 读取环境变量
	overrideFromEnv(&c)

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

func overrideFromEnv(c *config.Config) {
	c.Host = getEnv("HOST", c.Host)
	c.Port = int(getEnvInt("PORT", int64(c.Port)))
	c.Timeout = getEnvInt("TIMEOUT", c.Timeout)
	c.Log.Path = getEnv("LOG_PATH", c.Log.Path)
	c.MySqlDataSource = getEnv("MYSQL_DATASOURCE", c.MySqlDataSource)

	// 链路追踪配置

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return int64(intValue)
}
