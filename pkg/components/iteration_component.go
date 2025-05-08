package components

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/core"
)

type IterationComponent struct {
	config IterationConfig
}

type IterationConfig struct {
	IterationType  string              `json:"iterationType"`
	IterationValue core.Inputs         `json:"iterationValue"`
	SubWorkflowId  string              `json:"subWorkflowId"`
	workflowEngine core.WorkflowEngine `json:"-"`
}

var iterationComponentPool = sync.Pool{
	New: func() interface{} {
		return &IterationComponent{}
	},
}

func NewIterationComponent(e core.WorkflowEngine, config json.RawMessage) (*IterationComponent, error) {
	c := iterationComponentPool.Get().(*IterationComponent)
	var iteraConfig IterationConfig
	if err := sonic.Unmarshal(config, &iteraConfig); err != nil {
		return nil, errors.New("解析迭代执行组件配置失败: " + err.Error())
	}
	iteraConfig.workflowEngine = e
	c.config = iteraConfig
	return c, nil
}

const (
	iterationTotal = "iterationTotal"
)

// AnalyzeInputs implements Component.
func (i *IterationComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	execCtx := ctx.(*core.ExecutionContext)
	valMap, err := core.ParseNodeInputs([]core.Inputs{i.config.IterationValue}, execCtx)
	if err != nil {
		return nil, errors.New("解析迭代值失败: " + err.Error())
	}

	if len(valMap) == 0 {
		return nil, errors.New("迭代值为空")
	}

	// 获取第一个值
	var value any
	for _, v := range valMap {
		value = v
		break
	}

	input := map[string]any{}
	switch i.config.IterationType {
	case "array":
		arr, ok := value.([]any)
		if !ok {
			return nil, errors.New("迭代值必须是数组类型")
		}
		input[i.config.IterationValue.Name] = arr
	case "index":
		num, ok := value.(int64)
		if !ok {
			return nil, errors.New("迭代值必须是数字类型")
		}
		input[iterationTotal] = num
	default:
		return nil, errors.New("不支持的迭代类型: " + i.config.IterationType)
	}
	return input, nil
}

// Execute implements Component.
func (i *IterationComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	r := []any{}
	workflowID := i.config.SubWorkflowId
	switch i.config.IterationValue.Type[0] {
	case "array":
		{
			// 封装参数
			inputMap, ok := input.(map[string]any)
			if !ok {
				return nil, errors.New("输入类型不匹配")
			}
			iterVal, ok := inputMap[i.config.IterationValue.Name].([]any)
			if !ok {
				return nil, errors.New("迭代值类型不匹配")
			}
			// results := []core.NodeResult{}
			for idx, v := range iterVal {
				inputMap[i.config.IterationValue.Name] = v
				inputMap["_index"] = idx
				serialID := uuid.New().String()
				err := i.config.workflowEngine.ExecuteWorkflow(ctx, workflowID, serialID, inputMap)
				if err != nil {
					return nil, errors.New("迭代执行失败: " + err.Error())
				}
				result, ok := i.config.workflowEngine.GetNodeResult(workflowID, serialID, "end-item-node")
				if !ok {
					return nil, errors.New("获取迭代执行结果失败")
				}
				// results = append(results, *result)
				r = append(r, result.Output)
				logx.Debugf(" 迭代索引: %d\n 迭代参数: %+v\n 迭代结果: %+v\n", idx, inputMap, result)
			}
		}
	case "index":
		{
			// 封装参数
			inputMap, ok := input.(map[string]any)
			if !ok {
				return nil, errors.New("输入类型不匹配")
			}
			iterVal, ok := inputMap[iterationTotal].(int64)
			if !ok {
				return nil, errors.New("迭代值类型不匹配")
			}
			logx.Debugf("迭代内容: %d\n", iterVal)
		}
	}
	ret := make(map[string]any, 1)
	ret["result"] = r
	result := core.Result{
		Output: ret,
		Route:  []string{Success},
	}
	return &result, nil
}

// Validate implements Component.
func (i *IterationComponent) Validate() []core.ValidationError {
	return nil
}

// var _ Component = new(IterationComponent)
