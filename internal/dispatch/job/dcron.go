package job

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dcron-contrib/commons/dlog"
	"github.com/dcron-contrib/redisdriver"
	"github.com/libi/dcron"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
)

func InitDcron(ctx *svc.ServiceContext) {
	// 创建一个包装了logx的Logger实现
	logger := &dlog.StdLogger{
		Log:        &logxWrapper{},
		LogVerbose: false,
	}
	driver := redisdriver.NewDriver(ctx.RedisClient)

	dctx, cancel := context.WithCancel(context.Background())

	// 初始化dcron
	d := dcron.NewDcronWithOption(
		"workflow",
		driver,
		dcron.WithLogger(logger),
		dcron.WithHashReplicas(10),
		dcron.WithNodeUpdateDuration(time.Second*10),
	)

	DispatcherManager = &DcronManager{
		Dcron:  d,
		Ctx:    dctx,
		Cancel: cancel,
	}

	// 异步启动dcron
	go func() {
		d.Start()

		// 创建信号通道
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		// 监听上下文取消信号和系统停止信号
		go func() {
			select {
			case <-dctx.Done():
				d.Stop()
			case <-sigChan:
				d.Stop()
				cancel()
			}
		}()
		// Initialize system tasks
		// System tasks include:
		// 1. ProbDatasourceJob - Data source probe task
		// 2. SyncDatasourceJob - Data source sync task
		var systemJobCount int
		for _, jobConfig := range ctx.Config.Job {
			if jobConfig.Enable {
				var jobInstance dcron.Job
				switch jobConfig.Name {
				case ProbDatasourceJobName:
					jobInstance = &ProbDatasourceJob{}
					systemJobCount++
				case SyncDatasourceJobName:
					jobInstance = &SyncDatasourceJob{}
					systemJobCount++
				default:
					logx.Errorf("Unknown job name: %s", jobConfig.Name)
					continue
				}
				// If it's a 6-field expression, remove seconds field to make it 5-field
				cronExpr := jobConfig.Cron
				if len(strings.Fields(cronExpr)) == 6 {
					cronExpr = strings.Join(strings.Fields(cronExpr)[1:], " ")
				}
				if err := d.AddJob(jobConfig.Name, cronExpr, jobInstance); err != nil {
					logx.Errorf("Failed to add %s job: %v", jobConfig.Name, err)
				}
			}
		}

		// Initialize enabled tasks
		// 3. ChainJob - User-defined tasks loaded from database
		jobs, err := ctx.JobModel.FindByOn(context.Background())
		if err != nil {
			logx.Errorf("find job server error: %s\n", err.Error())
			return
		}
		var chainJobCount int
		for _, cjob := range jobs {
			jobInstance := &ChainJob{JobId: cjob.JobId, CanvasId: cjob.WorkspaceId}
			if err := d.AddJob(cjob.JobId, cjob.JobCron, jobInstance); err != nil {
				logx.Errorf("Failed to add %s job: %v", cjob.JobId, err)
			}
			chainJobCount++
		}
		logx.Infof("Dcron initialization completed. Loaded %d system jobs and %d chain jobs. Total jobs: %d",
			systemJobCount, chainJobCount, systemJobCount+chainJobCount)
	}()
}

// logx的包装器,实现PrintfLogger接口
type logxWrapper struct{}

func (l *logxWrapper) Printf(format string, args ...any) {
	logx.Infof(format, args...)
}
