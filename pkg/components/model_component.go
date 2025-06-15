package components

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/rotisserie/eris"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/chain"
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
	ModelId           int64                `json:"modelId"`
	Tools             []string             `json:"tools"`
	SystemPrompt      string               `json:"systemPrompt"`
	UserPrompt        string               `json:"userPrompt"`
	OutputType        string               `json:"outputType"`
	Retry             int64                `json:"retry"`
	Timeout           int64                `json:"timeout"`
	ErrorHandlingMode string               `json:"errorHandlingMode"`
	ExceptionConfig   ExceptionConfig      `json:"exceptionConfig"`
	NodeDataOutputs   core.NodeDataOutputs `json:"output"` // 输出定义
	OutputSchema      string               `json:"outputSchema"`
}

func (c *ModelComponent) Name() string {
	return Model
}

func NewModelComponent(config any) (*ModelComponent, error) {
	jsonConfig, err := sonic.Marshal(config)
	if err != nil {
		return nil, err
	}
	component := modelComponentPool.Get().(*ModelComponent)
	var modelConfig ModelConfig
	err = sonic.Unmarshal(jsonConfig, &modelConfig)
	if err != nil {
		return nil, err
	}
	modelConfig.ExceptionConfig = ExceptionConfig{
		Timeout:    modelConfig.Timeout,
		RetryTimes: int(modelConfig.Retry),
	}
	// 根据 NodeDataOutputs 输出定义，生成输出 json schema
	modelConfig.OutputSchema = generateOutputSchema(modelConfig.NodeDataOutputs)
	component.modelConfig = modelConfig
	return component, nil
}

func generateOutputSchema(outputs core.NodeDataOutputs) string {
	schema := map[string]interface{}{
		"$schema":    "http://json-schema.org/draft-07/schema#",
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   make([]string, 0),
	}

	for name, output := range outputs.Properties {
		propertySchema := map[string]interface{}{
			"type": convertTypeToJsonSchema(output.Type),
		}
		if output.Default != nil {
			propertySchema["default"] = output.Default
		}
		schema["properties"].(map[string]interface{})[name] = propertySchema
		if output.IsPropertyRequired {
			schema["required"] = append(schema["required"].([]string), name)
		}
	}

	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

// convertTypeToJsonSchema 将内部类型转换为 JSON Schema 类型
func convertTypeToJsonSchema(t string) string {
	switch strings.ToLower(t) {
	case "string":
		return "string"
	case "number", "float", "double":
		return "number"
	case "integer", "int":
		return "integer"
	case "boolean", "bool":
		return "boolean"
	case "array":
		return "array"
	case "object":
		return "object"
	case "null":
		return "null"
	default:
		return "string"
	}
}

func (c *ModelComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	var systemPrompt = c.modelConfig.SystemPrompt
	var userPrompt = c.modelConfig.UserPrompt

	inputMap, ok := input.(map[string]any)
	if !ok {
		return nil, eris.New("input type mismatch")
	}
	systemPromptAfter, err := c.replaceExprs(systemPrompt, inputMap)
	if err != nil {
		return nil, err
	}
	userPromptAfter, err := c.replaceExprs(userPrompt, inputMap)
	if err != nil {
		return nil, err
	}
	model, err := chain.AgentChain.NewChatModel(ctx, c.modelConfig.ModelId)
	if err != nil {
		logx.Errorf("new chat model failed: %s", err)
		return &core.Result{
			Output: map[string]any{},
			Route:  []string{Failed},
		}, err
	}
	// 输出结果类型
	outputType := c.modelConfig.OutputType
	if outputType == "" {
		outputType = "json"
	}
	resultPrompt := ""
	switch outputType {
	case "json":
		resultPrompt = "请按照以下jsonSchema的定义返回一个 json 对象,不要返回任何其他内容:" + c.modelConfig.OutputSchema
	case "string":
		resultPrompt = "请返回字符串格式结果,不要返回任何其他内容"
	}
	logx.Debugf("[MODEL] systemPrompt: %s\n%s\n", systemPromptAfter, resultPrompt)
	logx.Debugf("[MODEL] userPrompt: %s\n", userPromptAfter)

	// 调用模型
	result, err := model.Invoke(ctx, []*schema.Message{
		{
			Role:    "system",
			Content: systemPromptAfter + "\n" + resultPrompt,
		},
		{
			Role:    "user",
			Content: userPromptAfter,
		},
	})
	if err != nil {
		return &core.Result{
			Output: map[string]any{},
			Route:  []string{Failed},
		}, err
	}
	logx.Debugf("[MODEL] result: %s", result.Content)
	resultMap := map[string]any{}
	switch outputType {
	case "json":
		if result.Content != "" {
			err = json.Unmarshal([]byte(result.Content), &resultMap)
			if err != nil {
				logx.Errorf("[MODEL] unmarshal result failed: %s,content:%s", err, result.Content)
				// todo 再加一次大模型输出尝试
				result, err = model.Invoke(ctx, []*schema.Message{
					{
						Role:    "system",
						Content: result.Content + "\n" + resultPrompt,
					},
				})
				if err != nil {
					logx.Errorf("[MODEL] invoke model failed: %s", err)
					return &core.Result{
						Output: map[string]any{},
						Route:  []string{Failed},
					}, err
				}
				err = json.Unmarshal([]byte(result.Content), &resultMap)
				if err != nil {
					logx.Errorf("[MODEL] unmarshal result failed: %s,content:%s", err, result.Content)
					return &core.Result{
						Output: map[string]any{},
						Route:  []string{Failed},
					}, err
				}
			}
		}
	case "string":
		resultMap["content"] = result.Content
	}
	resultMap["completionTokens"] = result.ResponseMeta.Usage.CompletionTokens
	resultMap["promptTokens"] = result.ResponseMeta.Usage.PromptTokens
	resultMap["totalTokens"] = result.ResponseMeta.Usage.TotalTokens
	return &core.Result{
		Output: resultMap,
		Route:  []string{Success},
	}, nil
}

func (c *ModelComponent) replaceExprs(expr string, inputMap map[string]any) (string, error) {
	re := regexp.MustCompile(`{{.*?}}`)
	matches := re.FindAllString(expr, -1)

	var ok bool
	var value any
	for _, field := range matches {
		cutFiled := strings.Trim(field, "{{}}")
		if value, ok = inputMap[cutFiled]; !ok {
			// 本节点 input 里没有该变量
			return "", eris.Errorf("node input does not exist variable [%s]", cutFiled)
		}
		if value == nil {
			logx.Errorf("[MODEL] value is nil,cutFiled:%s", cutFiled)
			value = reflect.Zero(reflect.TypeOf(value)).Interface()
		}
		expr = strings.ReplaceAll(expr, field, fmt.Sprintf("%v", value))

	}
	return expr, nil
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
