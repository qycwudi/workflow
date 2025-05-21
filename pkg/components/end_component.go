package components

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/core"
)

// EndComponent 结束组件
type EndComponent struct {
	config EndConfig
}

type EndConfig struct {
	OutPutType string `json:"outPutType"` // 输出类型 json
}

var endComponentPool = sync.Pool{
	New: func() interface{} {
		return &EndComponent{}
	},
}

func NewEndComponent(config json.RawMessage) (*EndComponent, error) {
	var endConfig EndConfig
	if err := sonic.Unmarshal(config, &endConfig); err != nil {
		return nil, errors.New("解析结束组件配置失败: " + err.Error())
	}
	// 使用pool
	component := endComponentPool.Get().(*EndComponent)
	component.config = endConfig
	return component, nil
}

func (c *EndComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	logx.Debugf("结束组件执行: %+v\n", input)
	return &core.Result{
		Route:  []string{Success},
		Output: input,
	}, nil
}

func (c *EndComponent) Validate() []core.ValidationError {
	return nil
}

func (c *EndComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	return nil, nil
}

func (c *EndComponent) Exception() ExceptionConfig {
	return ExceptionConfig{}
}

func (c *EndComponent) Clear() {
	endComponentPool.Put(c)
}
