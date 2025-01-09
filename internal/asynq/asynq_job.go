package asynq

import (
	"log"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
)

// 创建Asynq客户端
func NewAsynqJob(ctx *svc.ServiceContext) {
	// 每个副本同一时间只运行一次
	singleProvider := &singleJobProvider{}
	for _, job := range ctx.Config.Job {
		// logx.Infof("create job: %s, enable: %t, cron: %s, topic: %s", job.Name, job.Enable, job.Cron, job.Topic)
		if job.Enable {
			// 解析cron表达式
			cronExpr := job.Cron
			// 尝试解析6段cron表达式,如果解析失败则使用原始表达式
			if _, err := cron.ParseStandard(cronExpr); err != nil {
				// 尝试将6段表达式转为5段,asynq 只支持 5 段
				parts := strings.Split(cronExpr, " ")
				if len(parts) == 6 {
					cronExpr = strings.Join(parts[1:], " ")
				}
			}

			// 如果是 * * * * * 则设置为每10分钟执行一次
			if cronExpr == "* * * * *" {
				cronExpr = "0/1 * * * *"
			}

			singleProvider.cycleCron = append(singleProvider.cycleCron, Config{
				Cron:  cronExpr,
				Topic: job.Topic,
			})

		}
	}

	mgr, err := asynq.NewPeriodicTaskManager(asynq.PeriodicTaskManagerOpts{
		RedisUniversalClient:       ctx.RedisClient,
		PeriodicTaskConfigProvider: singleProvider,
		SyncInterval:               1 * time.Hour, // 配置更新时间
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := mgr.Run(); err != nil {
			logx.Errorf("Asynq job run failed: %s", err.Error())
		}
	}()
}

type singleJobProvider struct {
	cycleCron []Config
}

type Config struct {
	Cron  string
	Topic string
}

// 更新配置
func (p *singleJobProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	var configs []*asynq.PeriodicTaskConfig
	for _, cfg := range p.cycleCron {
		logx.Infof("create job: %s, cron: %s, topic: %s", cfg.Topic, cfg.Cron, cfg.Topic)
		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: cfg.Cron, Task: asynq.NewTask(cfg.Topic, nil)})
	}
	return configs, nil
}
