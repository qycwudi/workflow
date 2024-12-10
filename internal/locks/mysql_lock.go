package locks

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
)

type MysqlLock struct {
	model model.LocksModel
}

func NewMysqlLock(model model.LocksModel) Lock {
	logx.Info("init mysql lock")
	return &MysqlLock{
		model: model,
	}
}

func (l *MysqlLock) Acquire(ctx context.Context, lockName string, ownerId string, timeout int) (bool, error) {
	return l.model.AcquireLock(ctx, lockName, ownerId, timeout)
}

func (l *MysqlLock) Release(ctx context.Context, lockName string, ownerId string) error {
	return l.model.ReleaseLock(ctx, lockName, ownerId)
}
