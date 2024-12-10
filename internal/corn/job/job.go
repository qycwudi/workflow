package job

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
)

// 定义任务执行函数类型
type JobFunc func(args ...interface{})

// 通用的任务调度函数
func RunScheduledTask(elector gocron.Elector, corn string, jobFunc JobFunc, args ...interface{}) error {
	s, err := gocron.NewScheduler(gocron.WithDistributedElector(elector))
	if err != nil {
		return fmt.Errorf("failed to create scheduler: %v", err)
	}
	j, err := s.NewJob(
		gocron.CronJob(corn, true),
		gocron.NewTask(jobFunc, args...),
	)
	if err != nil {
		return fmt.Errorf("failed to create job: %v", err)
	}

	fmt.Printf("%s job id: %s\n", args[0], j.ID())
	s.Start()

	// 等待程序退出信号
	<-context.Background().Done()

	if err = s.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown scheduler: %v", err)
	}
	fmt.Println("job shutdown")
	return nil
}
