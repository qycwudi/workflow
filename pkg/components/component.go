package components

import (
	"context"
	"errors"

	"github.com/bytedance/sonic"

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
	switch nodeType {
	case Start:
		return NewStartComponent()
	case End:
		return NewEndComponent()
	case HTTP:
		return NewHTTPComponent(inputs.Custom)
	case Code:
		return NewCodeComponent(inputs.Custom)
	case Model:
		return NewModelComponent(inputs.Custom)
	case Condition:
		jsonConfig, err := sonic.Marshal(inputs.Conditions)
		if err != nil {
			return nil, errors.New("component configuration serialization failed: " + err.Error())
		}
		return NewConditionComponent(jsonConfig)
	case Loop:
		return NewIterationComponent(e, inputs.BatchFor, inputs.Custom)
	case Database:
		return NewDatabaseComponent(inputs.Custom)
	}
	return nil, errors.New("Component type not found: " + nodeType)
}
