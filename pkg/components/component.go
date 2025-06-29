package components

import (
	"context"

	"workflow/pkg/core"
)

const (
	Success = "success"
	Failed  = "failed"

	Skip = "skip"

	True  = "true"
	False = "false"

	Else = "else"
)

const (
	Start     = "start"
	End       = "end"
	HTTP      = "http"
	Code      = "code"
	Condition = "condition"
	Model     = "llm"
	Loop      = "loop"
	StartItem = "start-item"
	EndItem   = "end-item"
	Database  = "sql"
)

// Component 定义组件核心接口
type Component interface {
	Validate() []core.ValidationError
	AnalyzeInputs(ctx context.Context) (any, error)
	Execute(ctx context.Context, input any) (*core.Result, error)
	Clear()
	Name() string
}

// ComponentFactory 组件工厂
func ComponentFactory(e core.WorkflowEngine, nodeType string, inputs core.NodeData) (Component, error) {
	return ComponentFactoryFixed(e, nodeType, inputs)
}
