package dispatch

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dcron-contrib/commons/dlog"
	"github.com/dcron-contrib/redisdriver"
	"github.com/libi/dcron"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/dispatch/job"
	"workflow/internal/svc"
)

type DcronManager struct {
	dcron  *dcron.Dcron
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

var DispatcherManager *DcronManager

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
		dcron:  d,
		ctx:    dctx,
		cancel: cancel,
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

		// 初始化系统任务
		for _, jobConfig := range ctx.Config.Job {
			if jobConfig.Enable {
				var jobInstance dcron.Job
				switch jobConfig.Name {
				case job.ProbDatasourceJobName:
					jobInstance = &job.ProbDatasourceJob{}
				case job.SyncDatasourceJobName:
					jobInstance = &job.SyncDatasourceJob{}
				case job.ChainJobName:
					jobInstance = &job.ChainJob{}
				default:
					logx.Errorf("Unknown job name: %s", jobConfig.Name)
					continue
				}
				// 如果是6位表达式,去掉秒位转成5位
				cronExpr := jobConfig.Cron
				if len(strings.Fields(cronExpr)) == 6 {
					cronExpr = strings.Join(strings.Fields(cronExpr)[1:], " ")
				}
				if err := d.AddJob(jobConfig.Name, cronExpr, jobInstance); err != nil {
					logx.Errorf("Failed to add %s job: %v", jobConfig.Name, err)
				}
			}
		}

		fmt.Println("dcron init success")
	}()
}

func (m *DcronManager) AddJob(name string, spec string, job dcron.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.dcron.AddJob(name, spec, job)
}

func (m *DcronManager) RemoveJob(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dcron.Remove(name)
}

func (m *DcronManager) Stop() {
	m.cancel()
}

// logx的包装器,实现PrintfLogger接口
type logxWrapper struct{}

func (l *logxWrapper) Printf(format string, args ...any) {
	logx.Infof(format, args...)
}
