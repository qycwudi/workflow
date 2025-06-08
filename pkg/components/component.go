package components

import (
	"context"
	"encoding/json"
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
	Iteration = "iteration"
	StartItem = "start-item"
	EndItem   = "end-item"
	Database  = "sql"
)

// Component 定义组件核心接口
type Component interface {
	Validate() []core.ValidationError
	AnalyzeInputs(ctx context.Context) (any, error)
	Execute(ctx context.Context, input any) (*core.Result, error)
	Exception() ExceptionConfig
	Clear()
}

// ComponentFactory 组件工厂
func ComponentFactory(e core.WorkflowEngine, nodeType string, inputs core.NodeData) (Component, error) {
	switch nodeType {
	case Start:
		return NewStartComponent()
	case End:
		return NewEndComponent(json.RawMessage("{}"))
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
	case Iteration:
		return NewIterationComponent(e, json.RawMessage("{}"))
	case StartItem:
		return NewStartItemComponent()
	case EndItem:
		return NewEndItemComponent(json.RawMessage("{}"))
	case Database:
		return NewDatabaseComponent(inputs.Custom)
	}
	return nil, errors.New("Component type not found: " + nodeType)
}

type ExceptionConfig struct {
	Timeout    int64 `json:"timeout"`
	RetryTimes int   `json:"retry_times"`
	// 处理方式，中断流程 / 返回设定内容 / 执行异常流程；如果是执行异常流程，会增加两个异常输出
	// 中断和执行异常流程 已经在 workflow 执行逻辑中实现了，组件只需要实现返回设定内
	// 发生异常时可以返回固定内容
	OutputOnError map[string]any `json:"output_on_error"`
}

func (e ExceptionConfig) Validate() []core.ValidationError {
	if e.RetryTimes > 10 || e.RetryTimes < 0 {
		return []core.ValidationError{
			{
				Message: "retry_times 值必须在 0 到 10 之间",
			},
		}
	}
	if e.Timeout <= 0 {
		return []core.ValidationError{
			{
				Message: "timeout 值必须大于 0",
			},
		}
	}
	return nil
}
