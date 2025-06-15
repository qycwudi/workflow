package engine

import (
	"sync"
	"time"
)

type WorkflowEngine struct {
	config *EngineConfig

	executorPool map[string]*Runnable // 工作流执行器池
	mu           sync.RWMutex         // 读写锁
}

// NewWorkflowEngine 创建新的工作流引擎
func NewWorkflowEngine(opts ...Option) *WorkflowEngine {
	config := DefaultConfig()
	engine := &WorkflowEngine{
		executorPool: make(map[string]*Runnable),
		config:       config,
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine
}

// Option 是工作流引擎的配置选项
type Option func(*WorkflowEngine)

// WithPoolSize 设置池大小
func WithPoolSize(size int) Option {
	return func(e *WorkflowEngine) {
		e.config.InitialPoolSize = size
	}
}

// WithMaxPoolSize 设置最大池大小
func WithMaxPoolSize(size int) Option {
	return func(e *WorkflowEngine) {
		e.config.MaxPoolSize = size
	}
}

// WithExecutionTimeout 设置执行超时
func WithExecutionTimeout(timeout time.Duration) Option {
	return func(e *WorkflowEngine) {
		e.config.ExecutionTimeout = timeout
	}
}

// WithMaxConcurrentWorkflows 设置最大并发工作流
func WithMaxConcurrentWorkflows(max int) Option {
	return func(e *WorkflowEngine) {
		e.config.MaxConcurrentWorkflows = max
	}
}
