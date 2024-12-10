package locks

import (
	"context"
)

var CustomLock Lock

// Lock 分布式锁接口
type Lock interface {
	// Acquire 获取锁
	// lockName: 锁的名称
	// ownerId: 锁持有者标识
	// timeout: 锁超时时间(秒)
	// 返回是否获取成功及错误信息
	Acquire(ctx context.Context, lockName string, ownerId string, timeout int) (bool, error)

	// Release 释放锁
	// lockName: 锁的名称
	// ownerId: 锁持有者标识
	// 返回错误信息
	Release(ctx context.Context, lockName string, ownerId string) error
}
