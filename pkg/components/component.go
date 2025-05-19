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
)

const (
	Start     = "start"
	End       = "end"
	HTTP      = "http"
	Codejs    = "codejs"
	Branch    = "branch"
	Model     = "model"
	Iteration = "iteration"
	StartItem = "start-item"
	EndItem   = "end-item"
	Database  = "database"
)

// Component 定义组件核心接口
type Component interface {
	Validate() []core.ValidationError
	AnalyzeInputs(ctx context.Context) (any, error)
	Execute(ctx context.Context, input any) (*core.Result, error)
	Clear()
}

// ComponentFactory 组件工厂
func ComponentFactory(e core.WorkflowEngine, nodeType string, nodeConfig *core.NodeDefinition) (Component, error) {
	jsonConfig, err := sonic.Marshal(nodeConfig.Config)
	if err != nil {
		return nil, errors.New("component configuration serialization failed: " + err.Error())
	}
	switch nodeType {
	case Start:
		return NewStartComponent()
	case End:
		return NewEndComponent(jsonConfig)
	case HTTP:
		return NewHTTPComponent(jsonConfig)
	case Codejs:
		return NewCodejsComponent(jsonConfig)
	case Model:
		return NewModelComponent(jsonConfig)
	case Branch:
		return NewBranchComponent(jsonConfig)
	case Iteration:
		return NewIterationComponent(e, jsonConfig)
	case StartItem:
		return NewStartItemComponent()
	case EndItem:
		return NewEndItemComponent(jsonConfig)
	case Database:
		return NewDatabaseComponent(jsonConfig)
	}
	return nil, errors.New("Component type not found: " + nodeType)
}
