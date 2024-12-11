package locks

import (
	"context"

	"workflow/internal/model"
)

type MysqlLock struct {
	model model.LocksModel
}

func (l *MysqlLock) Acquire(ctx context.Context, lockName string, ownerId string, timeout int) (bool, error) {
	return l.model.AcquireLock(ctx, lockName, ownerId, timeout)
}

func (l *MysqlLock) Release(ctx context.Context, lockName string, ownerId string) error {
	return l.model.ReleaseLock(ctx, lockName, ownerId)
}
