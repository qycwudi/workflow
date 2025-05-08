package components

import (
	"context"
	"sync"

	"workflow/pkg/core"
)

// StartItemComponent 迭代启动组件
type StartItemComponent struct {
}

var startItemComponentPool = sync.Pool{
	New: func() interface{} {
		return &StartItemComponent{}
	},
}

func NewStartItemComponent() (*StartItemComponent, error) {
	// 使用pool
	component := startItemComponentPool.Get().(*StartItemComponent)
	return component, nil
}

func (c *StartItemComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	return &core.Result{
		Route:  []string{Success},
		Output: input,
	}, nil
}

func (c *StartItemComponent) Validate() []core.ValidationError {
	return nil
}

func (c *StartItemComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	return nil, nil
}

func (c *StartItemComponent) Clear() {
	startItemComponentPool.Put(c)
}
