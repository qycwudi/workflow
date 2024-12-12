package servicecheck

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/corn/job"
	"workflow/internal/datasource"
	"workflow/internal/locks"
	"workflow/internal/model"
	"workflow/internal/svc"
)

// 探活数据源链接
var _ gocron.Elector = (*datasourceClientCheckElector)(nil)

type datasourceClientCheckElector struct{}

var datasourceClientCheckLockName = "datasource_client_probe_lock"

func (m *datasourceClientCheckElector) IsLeader(ctx context.Context) error {
	pid := strconv.Itoa(os.Getpid())
	has, err := locks.CustomLock.Acquire(ctx, datasourceClientCheckLockName, pid, 10)
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

func ProbeDatasourceClient(ctx *svc.ServiceContext, corn string) error {
	elector := &datasourceClientCheckElector{}
	ownerId := strconv.Itoa(os.Getpid())
	// Example task
	jobFunc := func(args ...interface{}) {
		context := context.Background()
		logx.Infof("%s datasource check start at: %s", args[0], time.Now().Format("2006-01-02 15:04:05"))
		total, datasourceList, err := ctx.DatasourceModel.FindDataSourcePageList(context, model.PageListBuilder{
			Switch: model.DatasourceSwitchOn,
		}, 1, 9999)
		if err != nil {
			logx.Error("Failed to get datasource list", err)
			// 释放锁
			if err = locks.CustomLock.Release(context, datasourceClientCheckLockName, ownerId); err != nil {
				logx.Errorf("Failed to release lock: %v", err)
			}
			return
		}

		logx.Infof("Successfully fetched %d datasources", total)
		successCount := 0
		failCount := 0
		for _, ds := range datasourceList {
			dsn := gjson.Get(ds.Config, "dsn").String()
			logx.Infof("Checking datasource: %s,%s", ds.Name, dsn)
			err := datasource.CheckDataSourceClient(ds.Type, dsn)
			nowStatus := model.DatasourceStatusConnected
			if err != nil {
				nowStatus = model.DatasourceStatusClosed
				failCount++
			} else {
				successCount++
			}
			if nowStatus != ds.Status {
				err = ctx.DatasourceModel.UpdateStatus(context, ds.Id, nowStatus)
				if err != nil {
					logx.Errorf("Datasource %d update status failed: %s", ds.Id, err.Error())
				}
			}
		}

		logx.Infof("Datasource check completed at: %s, success: %d, failed: %d",
			time.Now().Format("2006-01-02 15:04:05"), successCount, failCount)
		// 释放锁
		if err = locks.CustomLock.Release(context, datasourceClientCheckLockName, ownerId); err != nil {
			logx.Errorf("Failed to release lock: %v", err)
		}
	}

	return job.RunScheduledTask(elector, corn, jobFunc, "datasource_client_probe")
}
