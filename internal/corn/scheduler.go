package corn

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/config"
	"workflow/internal/corn/servicecheck"
	"workflow/internal/svc"
)

type Job struct{}

func NewJob(jobConfig config.JobConfig, ctx *svc.ServiceContext) {
	// Initialize scheduling task
	if jobConfig.DatasourceClientCheck.Enable {
		go func() {
			logx.Info("initialize servicecheck")
			err := servicecheck.Dispatch(ctx, jobConfig.DatasourceClientCheck.Cron)
			if err != nil {
				logx.Error("scheduling servicecheck task failed", err)
			}
		}()
	}

}

var _ gocron.Elector = (*defaultElector)(nil)

type defaultElector struct{}

func (m defaultElector) IsLeader(ctx context.Context) error {
	return nil
}
