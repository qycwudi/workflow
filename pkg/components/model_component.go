package components

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/bytedance/sonic"

	"workflow/pkg/core"
)

// ModelComponent 模型执行组件
type ModelComponent struct {
	modelConfig ModelConfig
}

var modelComponentPool = sync.Pool{
	New: func() interface{} {
		return &ModelComponent{}
	},
}

type ModelConfig struct {
	ModelId      int64          `json:"modelId"`
	Tools        []string       `json:"tools"`
	Params       map[string]any `json:"params"`
	SystemPrompt string         `json:"systemPrompt"`
	UserPrompt   string         `json:"userPrompt"`
}

func NewModelComponent(config json.RawMessage) (*ModelComponent, error) {
	component := modelComponentPool.Get().(*ModelComponent)
	var modelConfig ModelConfig
	err := sonic.Unmarshal(config, &modelConfig)
	if err != nil {
		return nil, err
	}
	component.modelConfig = modelConfig
	return component, nil
}

func (c *ModelComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	// execCtx := ctx.(*core.ExecutionContext)
	// var err error
	// var systemPrompt = c.modelConfig.SystemPrompt
	// var userPrompt = c.modelConfig.UserPrompt
	// modelParams := make(map[string]any)
	// if len(c.modelConfig.Params) != 0 {
	// 	modelParams, err = core.ParseNodeInputs(c.modelConfig.Params, execCtx)
	// 	if err != nil {
	// 		logx.Errorw("[模型组件] 解析参数失败",
	// 			logx.Field("错误", err))
	// 		return nil, err
	// 	}
	// 	// 替换 Prompt 中的 {{}}
	// 	for k, v := range modelParams {
	// 		logx.Debugf("[模型组件] 解析参数: %v", c.modelConfig.Params)
	// 		systemPrompt = strings.Replace(systemPrompt, "{{"+k+"}}", fmt.Sprintf("%v", v), -1)
	// 		userPrompt = strings.Replace(userPrompt, "{{"+k+"}}", fmt.Sprintf("%v", v), -1)
	// 	}

	// }
	// logx.Debugf("[模型组件] systemPrompt: %s", systemPrompt)
	// logx.Debugf("[模型组件] userPrompt: %s", userPrompt)

	// model, err := chain.AgentChain.NewChatModel(ctx, c.modelConfig.ModelId)
	// if err != nil {
	// 	logx.Errorf("new chat model failed: %s", err)
	// 	return &core.Result{
	// 		Output: map[string]any{},
	// 		Route:  []string{Failed},
	// 	}, err
	// }

	// result, err := model.Invoke(ctx, []*schema.Message{
	// 	{
	// 		Role:    "system",
	// 		Content: systemPrompt,
	// 	},
	// 	{
	// 		Role:    "user",
	// 		Content: userPrompt,
	// 	},
	// })
	// if err != nil {
	// 	return &core.Result{
	// 		Output: map[string]any{},
	// 		Route:  []string{Failed},
	// 	}, err
	// }
	// jsonResult, _ := sonic.Marshal(result)
	// logx.Debugf("[模型组件] 模型结果: %s", string(jsonResult))
	// mapResult := map[string]any{
	// 	"data": result.Content,
	// }
	// return &core.Result{
	// 	Output: mapResult,
	// 	Route:  []string{Success},
	// }, nil
	return nil, nil
}

func (c *ModelComponent) Validate() []core.ValidationError {
	return nil
}

func (c *ModelComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	return nil, nil
}

func (c *ModelComponent) Exception() ExceptionConfig {
	return ExceptionConfig{}
}

func (c *ModelComponent) Clear() {
	modelComponentPool.Put(c)
}
