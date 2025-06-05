package engine

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rotisserie/eris"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/components"
	"workflow/pkg/core"
)

// WorkflowEngine 工作流引擎 https://deepwiki.com/XXueTu/workflow/1-overview
type WorkflowEngine struct {
	config *EngineConfig // 引擎配置

	executorPool map[string]*Executor // 工作流执行器池
	mu           sync.RWMutex         // 读写锁

	cleanup    chan string   // 清理通道
	shutdownCh chan struct{} // 关闭信号通道

}

// NewWorkflowEngine 创建新的工作流引擎
func NewWorkflowEngine(opts ...Option) *WorkflowEngine {
	config := DefaultConfig()
	engine := &WorkflowEngine{
		executorPool: make(map[string]*Executor),
		config:       config,
		cleanup:      make(chan string, 100),
		shutdownCh:   make(chan struct{}),
	}

	for _, opt := range opts {
		opt(engine)
	}

	// 启动清理和指标收集协程
	engine.startCleanupRoutine()

	return engine
}

// createExecutor 创建执行器
func (e *WorkflowEngine) createExecutor(id string, def *core.WorkflowDef, plan *ExecutionPlan, condition map[string][]string) *Executor {
	e.mu.Lock()
	defer e.mu.Unlock()
	executor := &Executor{
		definition:      def,
		env:             make(map[string]any),
		executionPlan:   plan,
		conditionRouter: condition,
		totalNodes:      int64(len(def.Nodes)),
		status:          core.WorkflowStatusActive,
		createdAt:       time.Now(),
		defaultTTL:      e.config.DefaultContextTTL,
		shutdownCh:      make(chan struct{}),
	}

	if existing, exists := e.executorPool[id]; exists {
		atomic.StoreInt64(&existing.lastAccessed, time.Now().UnixNano())
		existing.status = core.WorkflowStatusDeploying
		// 复制已经存在的执行上下文,防止删除正在执行的上下文
		existing.execContexts.Range(func(key, value any) bool {
			executor.execContexts.Store(key, value)
			return true
		})
	}
	return executor
}

// DeregisterWorkflow 注销工作流
func (e *WorkflowEngine) DeregisterWorkflow(ctx context.Context, workflowID string, graceful bool) error {
	e.mu.Lock()
	executor, exists := e.executorPool[workflowID]
	if !exists {
		e.mu.Unlock()
		return eris.New("workflow not found: " + workflowID)
	}

	executor.status = core.WorkflowStatusShutdown
	e.mu.Unlock()

	if !graceful {
		return e.immediateDeregister(workflowID, executor)
	}

	return e.gracefulDeregister(ctx, workflowID, executor)
}

// immediateDeregister 立即注销工作流
func (e *WorkflowEngine) immediateDeregister(workflowID string, executor *Executor) error {
	e.mu.Lock()
	delete(e.executorPool, workflowID)
	e.mu.Unlock()
	close(executor.shutdownCh)
	return nil
}

// gracefulDeregister 优雅注销工作流
func (e *WorkflowEngine) gracefulDeregister(ctx context.Context, workflowID string, executor *Executor) error {
	go func() {
		timeout := time.After(30 * time.Second)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if atomic.LoadInt64(&executor.activeCount) == 0 {
					e.immediateDeregister(workflowID, executor)
					return
				}
			case <-timeout:
				e.immediateDeregister(workflowID, executor)
				return
			case <-ctx.Done():
				e.immediateDeregister(workflowID, executor)
				return
			}
		}
	}()

	return nil
}

// Clear 释放资源
func (e *WorkflowEngine) Clear(component components.Component) {
	component.Clear()
	logx.Debugf("[Workflow] Release component")
}

// 启动清理协程
func (e *WorkflowEngine) startCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(e.config.CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.cleanupExecContexts()
			case workflowID := <-e.cleanup:
				e.mu.Lock()
				delete(e.executorPool, workflowID)
				e.mu.Unlock()
			case <-e.shutdownCh:
				return
			}
		}
	}()
}

// cleanupExecContexts 清理执行上下文
func (e *WorkflowEngine) cleanupExecContexts() {
	now := time.Now()

	e.mu.RLock()
	executors := make([]*Executor, 0, len(e.executorPool))
	for _, exec := range e.executorPool {
		executors = append(executors, exec)
	}
	e.mu.RUnlock()

	for _, executor := range executors {
		// 清理每个执行器的上下文
		executor.execContexts.Range(func(key, value any) bool {
			id := key.(string)
			ctx, ok := value.(*core.ExecutionContext)
			if !ok {
				executor.execContexts.Delete(key)
				return true
			}

			// 检查是否已完成或过期
			if ctx.State.Status == core.StatusCompleted ||
				ctx.State.Status == core.StatusFailed ||
				now.After(ctx.Expiration) {
				executor.execContexts.Delete(id)
				core.ReleaseExecutionContext(ctx) // 释放上下文
			}
			return true
		})
	}
}

// Cleanup 释放资源
func (e *WorkflowEngine) Cleanup() {
	close(e.shutdownCh)

	// 关闭所有执行器
	e.mu.Lock()
	for id, executor := range e.executorPool {
		close(executor.shutdownCh)
		delete(e.executorPool, id)
		logx.Debugf("[Workflow] Close executor: %s", id)
	}
	e.mu.Unlock()
}

// GetExecutionContext 获取执行上下文
func (e *WorkflowEngine) GetExecutionContext(workflowID, serialID string) (*core.ExecutionContext, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	executor, ok := e.executorPool[workflowID]
	if !ok {
		return nil, eris.New("workflow not found: " + workflowID)
	}

	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return nil, eris.New("execution context not found: " + workflowID + ", " + serialID)
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return nil, eris.New("execution context type error: " + workflowID + ", " + serialID)
	}

	return execCtx, nil
}

func (e *WorkflowEngine) ClearExecutionContext(workflowID, serialID string) error {
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()
	if !ok {
		return eris.New("workflow not found: " + workflowID)
	}

	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return eris.New("execution context not found: " + workflowID + ", " + serialID)
	}

	execCtx, ok := val.(*core.ExecutionContext)
	executor.execContexts.Delete(serialID)
	core.ReleaseExecutionContext(execCtx)

	return nil
}

// ListWorkflows 列出所有工作流
func (e *WorkflowEngine) ListWorkflows() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	workflowIDs := make([]string, 0, len(e.executorPool))
	for id := range e.executorPool {
		workflowIDs = append(workflowIDs, id)
	}

	return workflowIDs
}

// GetWorkflowStatus 获取工作流状态
func (e *WorkflowEngine) GetWorkflowStatus(workflowID string) (core.WorkflowStatus, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	executor, ok := e.executorPool[workflowID]
	if !ok {
		return "", false
	}

	return executor.status, true
}

// PauseWorkflow 暂停工作流执行
func (e *WorkflowEngine) PauseWorkflow(ctx context.Context, workflowID string, serialID string) error {
	// 获取工作流执行器
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()

	if !ok {
		return eris.New("workflow not found: " + workflowID)
	}

	// 检查工作流状态
	if executor.status == core.WorkflowStatusShutdown {
		return eris.New("workflow is shutting down: " + workflowID)
	}

	// 使用sync.Map获取执行上下文
	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return eris.New("workflow execution instance not found: " + workflowID + ", " + serialID)
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return eris.New("workflow execution context type error: " + workflowID + ", " + serialID)
	}

	// 检查执行状态
	if execCtx.State.Status != core.StatusRunning {
		return eris.New("workflow execution status is not running: " + workflowID + ", " + serialID + ", current status: " + string(execCtx.State.Status))
	}

	// 更新执行状态为暂停
	execCtx.State.Status = core.StatusPaused

	// 获取并调用取消函数
	cancelVal, ok := execCtx.GetVariable("cancel")
	if !ok {
		return eris.New("failed to get cancel function: " + workflowID + ", " + serialID)
	}

	cancel, ok := cancelVal.(context.CancelFunc)
	if !ok {
		return eris.New("cancel function type error: " + workflowID + ", " + serialID)
	}

	// 调用取消函数
	cancel()
	return nil
}
