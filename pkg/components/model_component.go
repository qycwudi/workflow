package components

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"workflow/pkg/core"
)

// ModelComponent 模型执行组件
type ModelComponent struct {
	engine *GojaJsEngine
}

var modelComponentPool = sync.Pool{
	New: func() interface{} {
		return &ModelComponent{}
	},
}

type ModelConfig struct {
	Model        string        `json:"model"`
	Tools        []string      `json:"tools"`
	Params       []core.Inputs `json:"params"`
	SystemPrompt string        `json:"system_prompt"`
	UserPrompt   string        `json:"user_prompt"`
}

func NewModelComponent(config json.RawMessage) (*ModelComponent, error) {
	component := modelComponentPool.Get().(*ModelComponent)

	return component, nil
}

func (c *ModelComponent) Execute(ctx context.Context, input any) (*core.Result, error) {

	result, err := c.engine.Execute("main", input)
	if err != nil {
		return nil, errors.New("执行代码执行组件失败: " + err.Error())
	}

	return &core.Result{
		Output: result,
		Route:  []string{Success},
	}, nil
}

func (c *ModelComponent) Validate() []core.ValidationError {
	return nil
}

func (c *ModelComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	return nil, nil
}
