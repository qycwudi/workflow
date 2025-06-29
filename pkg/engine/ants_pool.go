package engine

import (
	"log"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DefaultAntsPoolSize = 100000           // 默认协程池大小，根据实际情况调整
	ExpiryDuration      = 10 * time.Second // worker空闲10秒后回收
)

var (
	globalAntsPool *ants.Pool
	once           sync.Once
)

// InitGlobalPool 初始化全局协程池
// 可以在应用启动时调用
func InitGlobalPool(opts ...ants.Option) error {
	var err error
	once.Do(func() {
		options := []ants.Option{
			ants.WithExpiryDuration(ExpiryDuration),
			ants.WithPanicHandler(func(err interface{}) {
				logx.Errorf("ants pool worker panic: %v", err)
			}),
			ants.WithNonblocking(false), // 如果池满，Submit会阻塞等待，推荐在高并发时设为false
			ants.WithPreAlloc(true),     // 预分配内存，减少运行时分配开销
		}
		options = append(options, opts...) // 允许外部传入更多自定义选项

		globalAntsPool, err = ants.NewPool(DefaultAntsPoolSize, options...)
		if err != nil {
			log.Fatalf("Failed to create ants pool: %v", err)
		}
		log.Printf("Ants pool initialized with size: %d, running: %d", globalAntsPool.Cap(), globalAntsPool.Running())
	})
	return err
}

// GetGlobalPool 获取全局协程池实例
// 如果未初始化，会panic，确保在使用前调用了InitGlobalPool
func GetGlobalPool() *ants.Pool {
	if globalAntsPool == nil {
		// 也可以在这里进行懒加载初始化，但通常推荐显式初始化
		log.Println("Warning: Global ants pool accessed before explicit initialization. Initializing with default settings.")
		if err := InitGlobalPool(); err != nil {
			// 如果懒加载初始化失败，这里应该panic或返回错误
			panic("Failed to lazy-initialize global ants pool: " + err.Error())
		}
	}
	return globalAntsPool
}

// ReleaseGlobalPool 释放全局协程池资源
// 可以在应用关闭时调用
func ReleaseGlobalPool() {
	if globalAntsPool != nil {
		globalAntsPool.Release()
		log.Println("Ants pool released.")
	}
}
