package corn

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/corn/servicecheck"
	"workflow/internal/svc"
)

type Job struct{}

func NewJob(ctx *svc.ServiceContext) {
	// Initialize scheduling task
	go func() {
		logx.Info("initialize servicecheck")
		err := servicecheck.Dispatch(ctx, "*/10 * * * * *")
		if err != nil {
			logx.Error("scheduling servicecheck task failed", err)
		}
	}()
}

var _ gocron.Elector = (*defaultElector)(nil)

type defaultElector struct{}

func (m defaultElector) IsLeader(ctx context.Context) error {
	return nil
}
