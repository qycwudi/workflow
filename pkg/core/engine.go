package core

import (
	"context"
)

type WorkflowEngine interface {
	ExecuteWorkflow(ctx context.Context, workflowID string, traceID string, params map[string]any, extra ContextExtra) (*ExecutionContextEnhanced, error)

	ListWorkflows() []string

	PauseWorkflow(ctx context.Context, workflowID string, traceID string) error
}

// WorkflowStatus 表示工作流的状态
type WorkflowStatus string

const (
	WorkflowStatusActive    WorkflowStatus = "active"
	WorkflowStatusDeploying WorkflowStatus = "deploying"
	WorkflowStatusShutdown  WorkflowStatus = "shutdown"
)
