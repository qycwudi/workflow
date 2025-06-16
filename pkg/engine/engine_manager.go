package engine

import (
	"context"

	"workflow/pkg/core"
)

// DeregisterWorkflow 注销工作流
func (e *WorkflowEngine) DeregisterWorkflow(ctx context.Context, workflowID string) error {
	e.mu.Lock()
	runnable, exists := e.runnablePool[workflowID]
	if !exists {
		e.mu.Unlock()
		return &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	runnable.status = core.WorkflowStatusShutdown
	e.mu.Unlock()

	return e.immediateDeregister(workflowID)
}

// immediateDeregister 立即注销工作流
func (e *WorkflowEngine) immediateDeregister(workflowID string) error {
	e.mu.Lock()
	delete(e.runnablePool, workflowID)
	e.mu.Unlock()
	return nil
}

// GetExecutionContext 获取执行上下文
func (e *WorkflowEngine) GetExecutionContext(workflowID, serialID string) (*core.ExecutionContext, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	runnable, ok := e.runnablePool[workflowID]
	if !ok {
		return nil, &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	val, ok := runnable.execContexts.Load(serialID)
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
	runnable, ok := e.runnablePool[workflowID]
	e.mu.RUnlock()
	if !ok {
		return &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	val, ok := runnable.execContexts.Load(serialID)
	if !ok {
		return &EngineManagerError{
			Message: "execution context not found: " + workflowID + ", " + serialID,
		}
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return &EngineManagerError{
			Message: "execution context type error: " + workflowID + ", " + serialID,
		}
	}
	runnable.execContexts.Delete(serialID)
	core.ReleaseExecutionContext(execCtx)

	return nil
}

// ListWorkflows 列出所有工作流
func (e *WorkflowEngine) ListWorkflows() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	workflowIDs := make([]string, 0, len(e.runnablePool))
	for id := range e.runnablePool {
		workflowIDs = append(workflowIDs, id)
	}

	return workflowIDs
}

// PauseWorkflow 暂停工作流执行
func (e *WorkflowEngine) PauseWorkflow(ctx context.Context, workflowID string, traceId string) error {
	// 获取工作流执行器
	e.mu.RLock()
	runnable, ok := e.runnablePool[workflowID]
	e.mu.RUnlock()

	if !ok {
		return &EngineManagerError{
			Message: "workflow not found: " + workflowID,
		}
	}

	// 检查工作流状态
	if runnable.status == core.WorkflowStatusShutdown {
		return &EngineManagerError{
			Message: "workflow is shutting down: " + workflowID,
		}
	}

	// 使用sync.Map获取执行上下文
	val, ok := runnable.execContexts.Load(traceId)
	if !ok {
		return &EngineManagerError{
			Message: "workflow execution instance not found: " + workflowID + ", " + traceId,
		}
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return &EngineManagerError{
			Message: "workflow execution context type error: " + workflowID + ", " + traceId,
		}
	}
	execCtx.Cancel()
	return nil
}
