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

// EndItemComponent 迭代结束组件
type EndItemComponent struct {
	config EndItemConfig
}

type EndItemConfig struct {
	OutPutType string `json:"outPutType"` // 输出类型 json
}

var endItemComponentPool = sync.Pool{
	New: func() interface{} {
		return &EndItemComponent{}
	},
}

func NewEndItemComponent(config json.RawMessage) (*EndItemComponent, error) {
	var endConfig EndItemConfig
	if err := sonic.Unmarshal(config, &endConfig); err != nil {
		return nil, errors.New("解析迭代结束组件配置失败: " + err.Error())
	}
	return &EndItemComponent{config: endConfig}, nil
}

func (c *EndItemComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	logx.Debugf("迭代结束组件执行: %+v\n", input)
	return &core.Result{
		Route:  []string{Success},
		Output: input,
	}, nil
}

func (c *EndItemComponent) Validate() []core.ValidationError {
	return nil
}

func (c *EndItemComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	return nil, nil
}
