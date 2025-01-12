package processor

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/datasource"
	"workflow/internal/enum"
	"workflow/internal/model"
)

const (
	TOPIC_DATA_SOURCE_CLIENT_PROBE = "datasource_client_probe"
)

type DatasourceClientProbeProcessor struct {
	datasourceModel model.DatasourceModel
}

type DatasourceClientProbePayload struct {
}

func (processor *DatasourceClientProbeProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	logx.Infof("%s start at: %s", TOPIC_DATA_SOURCE_CLIENT_PROBE, time.Now().Format("2006-01-02 15:04:05"))
	total, datasourceList, err := processor.datasourceModel.FindDataSourcePageList(ctx, model.PageListBuilder{
		Switch: model.DatasourceSwitchOn,
	}, 1, 9999)
	if err != nil && err != model.ErrNotFound {
		logx.Errorf("%s Failed to get datasource list: %v", TOPIC_DATA_SOURCE_CLIENT_PROBE, err)
		return err
	}

	logx.Infof("%s Successfully fetched %d datasources", TOPIC_DATA_SOURCE_CLIENT_PROBE, total)
	successCount := 0
	failCount := 0
	skipCount := 0
	for _, ds := range datasourceList {
		// 跳过fileServer
		if ds.Type == enum.FileServerType.String() {
			skipCount++
			logx.Infof("%s skip: %d, %s", TOPIC_DATA_SOURCE_CLIENT_PROBE, ds.Id, ds.Type)
			continue
		}
		err := datasource.CheckDataSourceClient(enum.DBType(ds.Type), ds.Config)
		nowStatus := model.DatasourceStatusConnected
		if err != nil {
			nowStatus = model.DatasourceStatusClosed
			failCount++
		} else {
			successCount++
		}
		if nowStatus != ds.Status {
			err = processor.datasourceModel.UpdateStatus(ctx, ds.Id, nowStatus)
			if err != nil {
				logx.Errorf("%s Datasource %d update status failed: %s", TOPIC_DATA_SOURCE_CLIENT_PROBE, ds.Id, err.Error())
			}
		}
	}
	logx.Infof("%s end at: %s, success: %d, failed: %d, skip: %d",
		TOPIC_DATA_SOURCE_CLIENT_PROBE, time.Now().Format("2006-01-02 15:04:05"), successCount, failCount, skipCount)
	return nil
}

func NewDatasourceClientProbeProcessor(datasourceModel model.DatasourceModel) *DatasourceClientProbeProcessor {
	return &DatasourceClientProbeProcessor{
		datasourceModel: datasourceModel,
	}
}
