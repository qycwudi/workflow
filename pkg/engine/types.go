package engine

import (
	"sync"
	"time"

	"workflow/pkg/components"
	"workflow/pkg/core"
)

type EngineConfig struct {
	InitialPoolSize        int // 初始池大小
	MaxPoolSize            int // 最大池大小
	MaxConcurrentWorkflows int // 最大并发工作流
}

type Runnable struct {
	definition      *core.Graph
	nodesNum        int64
	executionPlan   *ExecutionPlan
	conditionRouter map[string][]string
	execContexts    sync.Map            // 使用sync.Map替代map+mutex提高并发性能
	status          core.WorkflowStatus // 工作流状态
	createdAt       time.Time           // 创建时间
	defaultTTL      time.Duration       // 默认TTL
	shutdownCh      chan struct{}       // 关闭通道
}

// ExecutionTask 表示一个执行任务
type ExecutionTask struct {
	NodeID     string
	Component  components.Component
	Context    *core.ExecutionContext
	Engine     *WorkflowEngine
	Runnable   *Runnable
	WorkflowID string
	TraceID    string
	Step       int64
}

// ExecutionPlan 执行计划
type ExecutionPlan struct {
	Phases []ExecutionPhase
}

// ExecutionPhase 执行阶段
type ExecutionPhase struct {
	Nodes []*WorkflowNode
}

// WorkflowNode 工作流节点
type WorkflowNode struct {
	*core.Nodes
	Component components.Component
}

// DefaultConfig 返回默认配置
func DefaultConfig() *EngineConfig {
	return &EngineConfig{
		InitialPoolSize:        100,
		MaxPoolSize:            10000,
		MaxConcurrentWorkflows: 1000,
	}
}
