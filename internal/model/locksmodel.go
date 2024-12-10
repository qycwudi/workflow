package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LocksModel = (*customLocksModel)(nil)

type (
	// LocksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLocksModel.
	LocksModel interface {
		locksModel
		// 获取锁
		AcquireLock(ctx context.Context, lockName string, ownerId string, timeout int) (bool, error)
		// 释放锁
		ReleaseLock(ctx context.Context, lockName string, ownerId string) error
	}

	customLocksModel struct {
		*defaultLocksModel
	}
)

func NewLocksModel(conn sqlx.SqlConn) LocksModel {
	return &customLocksModel{
		defaultLocksModel: newLocksModel(conn),
	}
}

// AcquireLock implements LocksModel.
func (c *customLocksModel) AcquireLock(ctx context.Context, lockName string, ownerId string, timeout int) (bool, error) {
	var acquired bool
	err := c.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var locks Locks
		query := fmt.Sprintf("SELECT %s FROM %s WHERE lock_name = ? FOR UPDATE", locksRows, c.table)
		err := session.QueryRowCtx(ctx, &locks, query, lockName)
		if err == sqlc.ErrNotFound {
			// 如果不存在锁，则直接创建一条锁记录并标记为加锁
			_, err = session.ExecCtx(ctx, fmt.Sprintf("INSERT INTO %s (lock_name, is_locked, held_by, locked_time, timeout, updated_time) VALUES (?, 1, ?, NOW(), ?, NOW())", c.table),
				lockName, ownerId, timeout)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrLockAcquireFail, err)
			}
			acquired = true
			return nil
		} else if err == nil {
			// 如果锁记录存在，需要判断是否超时
			elapsed := time.Since(locks.LockedTime).Seconds()
			if locks.IsLocked == 1 || elapsed < float64(locks.Timeout) {
				// 锁被其他持有者持有，或未超时，返回获取失败
				acquired = false
				return nil
			}

			// 如果锁已超时或未锁定，则重新加锁
			_, err = session.ExecCtx(ctx, fmt.Sprintf("UPDATE %s SET is_locked = 1, held_by = ?, locked_time = NOW(), timeout = ?, updated_time = NOW() WHERE lock_name = ?", c.table),
				ownerId, timeout, lockName)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrLockAcquireFail, err)
			}
			acquired = true
			return nil
		}

		return fmt.Errorf("%w: %v", ErrLockNotFound, err)
	})

	if err != nil {
		return false, err
	}

	return acquired, nil
}

// ReleaseLock implements LocksModel.
func (c *customLocksModel) ReleaseLock(ctx context.Context, lockName string, ownerId string) error {
	query := fmt.Sprintf("UPDATE %s SET is_locked = 0, held_by = '', updated_time = NOW() WHERE lock_name = ? AND held_by = ?", c.table)
	result, err := c.conn.ExecCtx(ctx, query, lockName, ownerId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLockReleaseFail, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLockReleaseFail, err)
	}
	if rowsAffected == 0 {
		return ErrInvalidLockOwner
	}

	return nil
}

// LockError 锁相关错误
type LockError struct {
	Message string
}

func (e *LockError) Error() string {
	return e.Message
}

// 定义常见错误
var (
	ErrLockNotFound     = &LockError{Message: "lock not found"}
	ErrLockAcquireFail  = &LockError{Message: "failed to acquire lock"}
	ErrLockReleaseFail  = &LockError{Message: "failed to release lock"}
	ErrLockTimeout      = &LockError{Message: "lock timeout"}
	ErrInvalidLockOwner = &LockError{Message: "invalid lock owner"}
)
