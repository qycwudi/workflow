package servicecheck

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/corn/job"
	"workflow/internal/locks"
	"workflow/internal/model"
	"workflow/internal/svc"
)

// 探活数据源链接
var _ gocron.Elector = (*datasourceClientCheckElector)(nil)

type datasourceClientCheckElector struct{}

var datasourceClientCheckLockName = "datasource_client_check_lock"

func (m *datasourceClientCheckElector) IsLeader(ctx context.Context) error {
	has, err := locks.CustomLock.Acquire(ctx, datasourceClientCheckLockName, "system", 10)
	// logx.Infof("has: %v, err: %v", has, err)
	if err != nil {
		logx.Errorf(err.Error())
		return err
	}
	if !has {
		return fmt.Errorf("not leader")
	}
	return nil
}

func Dispatch(ctx *svc.ServiceContext, corn string) error {
	elector := &datasourceClientCheckElector{}
	// Example task
	jobFunc := func(args ...interface{}) {
		context := context.Background()
		logx.Infof("%s datasource check start at: %s", args[0], time.Now().Format("2006-01-02 15:04:05"))
		total, datasourceList, err := ctx.DatasourceModel.FindDataSourcePageList(context, model.PageListBuilder{
			Type:   "mysql",
			Status: "enable",
			Switch: 1,
		}, 1, 100)
		if err != nil {
			logx.Error("Failed to get datasource list", err)
			// 释放锁
			if err = locks.CustomLock.Release(context, datasourceClientCheckLockName, "system"); err != nil {
				logx.Errorf("Failed to release lock: %v", err)
			}
			return
		}

		logx.Infof("Successfully fetched %d datasources", total)
		for _, ds := range datasourceList {
			logx.Infof("Checking datasource: %d", ds.Id)
		}

		logx.Info("Datasource check completed at: ", time.Now().Format("2006-01-02 15:04:05"))
		// 释放锁
		if err = locks.CustomLock.Release(context, "datasource_client_check_lock", "system"); err != nil {
			logx.Errorf("Failed to release lock: %v", err)
		}
	}

	return job.RunScheduledTask(elector, corn, jobFunc, "mysql")
}
