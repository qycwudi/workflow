package servicecheck

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/corn/job"
	"workflow/internal/datasource"
	"workflow/internal/enum"
	"workflow/internal/model"
	"workflow/internal/svc"
)

// 探活数据源链接
var _ gocron.Elector = (*datasourceClientUpdateElector)(nil)

type datasourceClientUpdateElector struct{}

func (m *datasourceClientUpdateElector) IsLeader(ctx context.Context) error {
	return nil
}

func UpdateDatasourceClient(ctx *svc.ServiceContext, corn string) error {
	elector := &datasourceClientUpdateElector{}
	jobFunc := func(args ...interface{}) {
		context := context.Background()
		logx.Infof("%s datasource update start at: %s", args[0], time.Now().Format("2006-01-02 15:04:05"))
		total, datasourceList, err := ctx.DatasourceModel.FindDataSourcePageList(context, model.PageListBuilder{
			Switch: model.DatasourceSwitchOn,
		}, 1, 9999)
		if err != nil {
			logx.Error("Failed to get datasource list", err)
			return
		}

		logx.Infof("Successfully fetched %d datasources", total)
		successCount := 0
		failCount := 0
		skipCount := 0
		for _, ds := range datasourceList {
			// 跳过fileServer
			if ds.Type == enum.FileServerType.String() {
				skipCount++
				logx.Infof("datasource update skip: %d, %s", ds.Id, ds.Type)
				continue
			}
			err := datasource.DataSourcePool.UpdateDataSource(ds.Id, ds.Config, ds.Type, ds.Hash)
			if err != nil {
				logx.Errorf("Datasource %d update failed: %s", ds.Id, err.Error())
				failCount++
			} else {
				successCount++
			}
		}
		// 清理连接池中已删除的数据源连接
		poolIds := make(map[int64]bool)
		for _, ds := range datasourceList {
			poolIds[ds.Id] = true
		}

		// 遍历连接池,清理不存在于数据库的连接
		clearCount := 0
		for id := range datasource.DataSourcePool.GetHash() {
			if !poolIds[id] {
				if err := datasource.DataSourcePool.ClearDataSource(id); err != nil {
					logx.Errorf("clear datasource: %d, err: %v", id, err)
				} else {
					clearCount++
					logx.Infof("clear datasource: %d, success", id)
				}
			}
		}

		logx.Infof("Datasource update completed at: %s, total: %d, success: %d, failed: %d, cleared: %d, skip: %d",
			time.Now().Format("2006-01-02 15:04:05"), total, successCount, failCount, clearCount, skipCount)
	}

	return job.RunScheduledTask(elector, corn, jobFunc, "datasource_client_update")
}
