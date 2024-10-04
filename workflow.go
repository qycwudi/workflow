package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"os"
	"time"
	"workflow/internal/config"
	"workflow/internal/handler"
	"workflow/internal/rolego"
	"workflow/internal/svc"
	"workflow/internal/utils"
)

var configFile = flag.String("f", "etc/workflow-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// 读取配置文件
	conf.MustLoad(*configFile, &c)

	// 读取环境变量
	overrideFromEnv(&c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 自定义日志输出
	if c.OpenOB.Path != "" {
		httpWriter := utils.NewHTTPLogWriter(c.OpenOB.Path, c.OpenOB.UserName, c.OpenOB.Password, 5*time.Second)
		logx.AddWriter(logx.NewWriter(httpWriter))
	}

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// 注册链服务
	rolego.InitRoleServer()
	// 注册规则链
	rolego.InitRoleChain()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

/*
	LOG_MODE=json
	LOG_PATH=/app/logs

	OPENOB_PATH=http://10.99.43.100:5080/api/default/workflow/_json;
	OPENOB_USERNAME=root@123456789.com;
	OPENOB_PASSWORD=123456789;

	TELEMETRY_AUTHORIZATION=Basic cm9vdEAxMjM0NTY3ODkuY29tOkNRY1Y0WHFpaWtVR255TEw\\=;
	TELEMETRY_ENDPOINT=10.99.43.100:5081
*/

func overrideFromEnv(c *config.Config) {
	c.Log.Mode = getEnv("LOG_MODE", c.Log.Mode)
	c.Log.Path = getEnv("LOG_PATH", c.Log.Path)
	c.MySqlDataSource = getEnv("MYSQL_DATASOURCE", c.MySqlDataSource)

	// 系统监控配置
	c.Telemetry.Endpoint = getEnv("TELEMETRY_ENDPOINT", "")
	c.Telemetry.OtlpHeaders = map[string]string{
		"Authorization": getEnv("TELEMETRY_AUTHORIZATION", ""),
		"organization":  getEnv("TELEMETRY_ORGANIZATION", "default"),
		"stream-name":   getEnv("TELEMETRY_STREAM_NAME", "default"),
	}

	// 自定义追踪
	c.OpenOB.Path = getEnv("OPENOB_PATH", "")
	c.OpenOB.UserName = getEnv("OPENOB_USERNAME", "")
	c.OpenOB.Password = getEnv("OPENOB_PASSWORD", "")

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
