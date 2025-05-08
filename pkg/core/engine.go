package core

import (
	"context"
)

type WorkflowEngine interface {
	ExecuteWorkflow(ctx context.Context, workflowID string, serialID string, params map[string]any) error
	GetNodeResult(workflowID, serialID, nodeID string) (*NodeResult, bool)

	ListWorkflows() []string
	GetWorkflowStatus(workflowID string) (WorkflowStatus, bool)

	PauseWorkflow(ctx context.Context, workflowID string, serialID string) error
}

// WorkflowStatus 表示工作流的状态
type WorkflowStatus string

const (
	WorkflowStatusActive    WorkflowStatus = "active"
	WorkflowStatusDeploying WorkflowStatus = "deploying"
	WorkflowStatusShutdown  WorkflowStatus = "shutdown"
)
