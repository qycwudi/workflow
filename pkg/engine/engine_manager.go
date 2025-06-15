package engine

import (
	"context"

	"workflow/pkg/core"
)

// DeregisterWorkflow 注销工作流
func (e *WorkflowEngine) DeregisterWorkflow(ctx context.Context, workflowID string) error {
	e.mu.Lock()
	executor, exists := e.executorPool[workflowID]
	if !exists {
		e.mu.Unlock()
		return &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	executor.status = core.WorkflowStatusShutdown
	e.mu.Unlock()

	return e.immediateDeregister(workflowID)
}

// immediateDeregister 立即注销工作流
func (e *WorkflowEngine) immediateDeregister(workflowID string) error {
	e.mu.Lock()
	delete(e.executorPool, workflowID)
	e.mu.Unlock()
	return nil
}

// GetExecutionContext 获取执行上下文
func (e *WorkflowEngine) GetExecutionContext(workflowID, serialID string) (*core.ExecutionContext, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	executor, ok := e.executorPool[workflowID]
	if !ok {
		return nil, &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return nil, &EngineManagerError{
			Message: "execution context not found: " + workflowID + ", " + serialID,
		}
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return nil, &EngineManagerError{
			Message: "execution context type error: " + workflowID + ", " + serialID,
		}
	}

	return execCtx, nil
}

// ClearExecutionContext 清除执行上下文
func (e *WorkflowEngine) ClearExecutionContext(workflowID, serialID string) error {
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()
	if !ok {
		return &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return &EngineManagerError{
			Message: "execution context not found: " + workflowID + ", " + serialID,
		}
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
		return &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	// 检查工作流状态
	if executor.status == core.WorkflowStatusShutdown {
		return &EngineManagerError{
			Message: "workflow is shutting down: " + workflowID,
		}
	}

	// 使用sync.Map获取执行上下文
	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return &EngineManagerError{
			Message: "workflow execution instance not found: " + workflowID + ", " + serialID,
		}
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return &EngineManagerError{
			Message: "workflow execution context type error: " + workflowID + ", " + serialID,
		}
	}

	// 检查执行状态
	if execCtx.State.Status != core.StatusRunning {
		return &EngineManagerError{
			Message: "workflow execution status is not running: " + workflowID + ", " + serialID + ", current status: " + string(execCtx.State.Status),
		}
	}

	// 更新执行状态为暂停
	execCtx.State.Status = core.StatusPaused

	// 获取并调用取消函数
	cancelVal, ok := execCtx.GetVariable("cancel")
	if !ok {
		return &EngineManagerError{
			Message: "failed to get cancel function: " + workflowID + ", " + serialID,
		}
	}

	cancel, ok := cancelVal.(context.CancelFunc)
	if !ok {
		return &EngineManagerError{
			Message: "cancel function type error: " + workflowID + ", " + serialID,
		}
	}

	// 调用取消函数
	cancel()
	return nil
}
