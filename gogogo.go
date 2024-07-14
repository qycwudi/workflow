package main

import (
	"flag"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gogogo/internal/config"
	"gogogo/internal/handler"
	"gogogo/internal/svc"
	"log"
	"net/http"
)

var configFile = flag.String("f", "etc/gogogo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	// -// 开启监听
	// -consumer := mq.NewRockerConsumer(svc.GetOCRConsumer(ctx), svc.GetLLMConsumer(ctx))
	// -go func() {
	// 	-// 开启OCR监听
	// 	-consumer.ReceiveOCRMessage()
	// -}()
	// -go func() {
	// 	-// 开启LLM监听
	// 	-consumer.ReceiveLLMMessage()
	// -}()
	//
	// -// 关闭rockerMQ生产者
	// -defer func() {
	// 	-mq.CloseProducer(ctx.OcrProducer)
	// 	-mq.CloseProducer(ctx.LlmProducer)
	// -}()

	// 关闭asynq任务队列

	defer ctx.AsynqTaskClient.Close()

	// 开启asynq监控
	go func() {
		h := asynqmon.New(asynqmon.Options{
			RootPath: "/asynq",
			RedisConnOpt: asynq.RedisClientOpt{
				Addr:     c.RedisConfig.RedisAddr,
				DB:       c.RedisConfig.RedisDb,
				Username: c.RedisConfig.RedisUserName,
				Password: c.RedisConfig.RedisPassword,
			},
		})

		http.Handle(h.RootPath()+"/", h)
		println("Starting asynq monitor at 0.0.0.0:8889...")
		log.Fatal(http.ListenAndServe(":8889", nil))
	}()
	server.Start()
}
