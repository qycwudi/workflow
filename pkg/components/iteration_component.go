package components

import (
	"context"
	"errors"
	"sync"

	"github.com/bytedance/sonic"

	"workflow/pkg/core"
)

type IterationComponent struct {
	config IterationConfig
}

type IterationConfig struct {
	EndId string `json:"endId"`
	// IterationType  string              `json:"iterationType"`
	workflowEngine core.WorkflowEngine       `json:"-"`
	BatchFor       core.NodeDataInputsValues `json:"batchFor"`
}

var iterationComponentPool = sync.Pool{
	New: func() interface{} {
		return &IterationComponent{}
	},
}

func NewIterationComponent(e core.WorkflowEngine, batchFor core.NodeDataInputsValues, config any) (*IterationComponent, error) {
	jsonConfig, err := sonic.Marshal(config)
	if err != nil {
		return nil, errors.New("loop component configuration serialization failed: " + err.Error())
	}
	// 使用pool
	component := iterationComponentPool.Get().(*IterationComponent)
	var iteraConfig IterationConfig
	if err := sonic.Unmarshal(jsonConfig, &iteraConfig); err != nil {
		return nil, errors.New("loop component configuration serialization failed: " + err.Error())
	}
	iteraConfig.workflowEngine = e
	iteraConfig.BatchFor = batchFor
	component.config = iteraConfig
	return component, nil
}

// AnalyzeInputs implements Component.
func (i *IterationComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	// execCtx := ctx.(*core.ExecutionContext)
	// // 构造
	// inputs := map[string]core.NodeDataInputsValues{"batchFor": i.config.BatchFor}
	// inputValues := core.NodeDataInputs{
	// 	Properties: map[string]core.Properties{
	// 		"batchFor": {
	// 			Type: i.config.BatchFor.Type,
	// 		},
	// 	},
	// 	Required: []string{"batchFor"},
	// }
	// valMap, err := core.ParseNodeInputs(execCtx, inputs, inputValues)
	// if err != nil {
	// 	return nil, errors.New("解析迭代值失败: " + err.Error())
	// }

	// if len(valMap) == 0 {
	// 	return nil, errors.New("迭代值为空")
	// }

	// // 获取第一个值
	// var value any
	// for _, v := range valMap {
	// 	value = v
	// 	break
	// }

	// input := map[string]any{}
	// switch i.config.IterationType {
	// case "array":
	// 	arr, ok := value.([]any)
	// 	if !ok {
	// 		return nil, errors.New("迭代值必须是数组类型")
	// 	}
	// 	input[i.config.IterationValue.Name] = arr
	// case "index":
	// 	num, ok := value.(int64)
	// 	if !ok {
	// 		return nil, errors.New("迭代值必须是数字类型")
	// 	}
	// 	input[iterationTotal] = num
	// default:
	// 	return nil, errors.New("不支持的迭代类型: " + i.config.IterationType)
	// }
	// return input, nil
	return nil, nil
}

// Execute implements Component.
func (i *IterationComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	// r := []any{}
	// workflowID := i.config.SubWorkflowId
	// switch i.config.IterationValue.Type[0] {
	// case "array":
	// 	{
	// 		// 封装参数
	// 		inputMap, ok := input.(map[string]any)
	// 		if !ok {
	// 			return nil, errors.New("输入类型不匹配")
	// 		}
	// 		iterVal, ok := inputMap[i.config.IterationValue.Name].([]any)
	// 		if !ok {
	// 			return nil, errors.New("迭代值类型不匹配")
	// 		}
	// 		// results := []core.NodeResult{}
	// 		for idx, v := range iterVal {
	// 			inputMap[i.config.IterationValue.Name] = v
	// 			inputMap["_index"] = idx
	// 			serialID := uuid.New().String()
	// 			err := i.config.workflowEngine.ExecuteWorkflow(ctx, workflowID, serialID, inputMap)
	// 			if err != nil {
	// 				return nil, errors.New("迭代执行失败: " + err.Error())
	// 			}
	// 			result, ok := i.config.workflowEngine.GetNodeResult(workflowID, serialID, i.config.EndId)
	// 			if !ok {
	// 				return nil, errors.New("获取迭代执行结果失败")
	// 			}
	// 			// results = append(results, *result)
	// 			r = append(r, result.Output)
	// 			logx.Debugf(" 迭代索引: %d\n 迭代参数: %+v\n 迭代结果: %+v\n", idx, inputMap, result)
	// 		}
	// 	}
	// case "index":
	// 	{
	// 		// 封装参数
	// 		inputMap, ok := input.(map[string]any)
	// 		if !ok {
	// 			return nil, errors.New("输入类型不匹配")
	// 		}
	// 		iterVal, ok := inputMap[iterationTotal].(int64)
	// 		if !ok {
	// 			return nil, errors.New("迭代值类型不匹配")
	// 		}
	// 		logx.Debugf("迭代内容: %d\n", iterVal)
	// 	}
	// }
	// ret := make(map[string]any, 1)
	// ret["result"] = r
	// result := core.Result{
	// 	Output: ret,
	// 	Route:  []string{Success},
	// }
	// return &result, nil
	return nil, nil
}

// Validate implements Component.
func (i *IterationComponent) Validate() []core.ValidationError {
	return nil
}

func (c *IterationComponent) Exception() ExceptionConfig {
	return ExceptionConfig{}
}

// var _ Component = new(IterationComponent)

func (i *IterationComponent) Clear() {
	iterationComponentPool.Put(i)
}
