package components

import (
	"context"
	"errors"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/core"
)

type IterationComponent struct {
	config IterationConfig
}

type IterationConfig struct {
	NodeID         string              `json:"nodeId"`
	BatchForType   string              `json:"batchForType"`
	workflowEngine core.WorkflowEngine `json:"-"`

	BatchFor core.NodeDataInputsValues  `json:"batchFor"`
	Outputs  map[string]core.Properties `json:"outputs"`
}

var iterationComponentPool = sync.Pool{
	New: func() interface{} {
		return &IterationComponent{}
	},
}

const (
	itemKey = "items"
)

func (c *IterationComponent) Name() string {
	return Loop
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
	logx.Debugf("loop component configuration: %+v", iteraConfig)
	return component, nil
}

// AnalyzeInputs implements Component.
func (i *IterationComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	execCtx := ctx.(*core.ExecutionContextEnhanced)
	// 构造
	inputs := map[string]core.NodeDataInputsValues{"batchFor": i.config.BatchFor}
	inputValues := core.NodeDataInputs{
		Properties: map[string]core.Properties{
			"batchFor": {
				Type: "array",
				Item: core.NodeDataOutputs{
					Type: i.config.BatchForType,
				},
			},
		},
		Required: []string{"batchFor"},
	}
	valMap, err := core.ParseNodeInputs(execCtx, inputs, inputValues)
	if err != nil {
		return nil, errors.New("[loop] analyze inputs failed: " + err.Error())
	}

	if len(valMap) == 0 {
		return nil, errors.New("[loop] analyze inputs failed: input is empty")
	}

	// 获取第一个值
	var value any
	for _, v := range valMap {
		value = v
		break
	}

	input := map[string]any{}
	arr, ok := value.([]any)
	if !ok {
		return nil, errors.New("[loop] analyze inputs failed: input is not array")
	}
	input[itemKey] = arr

	return input, nil
}

// Execute implements Component.
func (i *IterationComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	execCtx, ok := ctx.(*core.ExecutionContextEnhanced)
	if !ok {
		return nil, errors.New("loop component context type error")
	}
	r := make(map[string]any, 0)
	/*
		{
			"result_1": [
				{
					"key": "value"
				}
			],
			"result_2": [
				{
					"key": "value"
				}
			]
		}
	*/
	workflowID := i.config.NodeID
	// 封装参数
	inputMap, ok := input.(map[string]any)
	if !ok {
		return nil, errors.New("[loop] execute failed: input type mismatch")
	}
	iterVal, ok := inputMap[itemKey].([]any)
	if !ok {
		return nil, errors.New("[loop] execute failed: item type mismatch")
	}
	for idx, v := range iterVal {
		// 迭代参数
		inputMap["item"] = v
		// 迭代索引
		inputMap["index"] = idx
		// 设置loop输出参数,用于子流程获取
		execCtx.SetVariable(workflowID+"_locals"+".output", inputMap)
		logx.Debugf("[loop] setVariable success: index: %d, output: %+v", idx, inputMap)
		subExecCtx, err := i.config.workflowEngine.ExecuteWorkflow(execCtx, workflowID, execCtx.TraceId, inputMap, core.ContextExtra{
			IsSub:             true,
			ParentWorkspaceId: execCtx.WorkspaceId,
			Index:             int64(idx),
			NodeNum:           int64(len(iterVal)),
		})
		if err != nil {
			// todo 错误处理机制
			return nil, errors.New("[loop] execute failed: " + err.Error())
		}
		// 获取指定节点输出
		subResult := make(map[string]any, 0)
		for _, output := range i.config.Outputs {
			key := output.Extra.InputKey
			outputKey := output.Extra.OutputKey
			subNodeId := output.Extra.NodeId
			outputValue, ok := subExecCtx.GetVariable(subNodeId + ".output")
			if !ok {
				return nil, errors.New("[loop] execute failed: get node result failed")
			}
			// 赋值
			subResult[key] = outputValue[outputKey]
		}
		// 把subResult append 到 r
		for key, value := range subResult {
			if _, ok := r[key]; !ok {
				r[key] = make([]any, 0)
			}
			// 类型断言并追加
			if arr, ok := r[key].([]any); ok {
				r[key] = append(arr, value)
			}
		}
		logx.Debugf("[loop] execute success: index: %d, item: %+v, result: %+v", idx, inputMap, r)
	}

	// 迭代完成后清空 sub_index，恢复到非子流程状态
	execCtx.Extra.Index = -1
	execCtx.Extra.IsSub = false

	result := core.Result{
		Output: r,
		Route:  []string{Success},
	}
	return &result, nil
}

// Validate implements Component.
func (i *IterationComponent) Validate() []core.ValidationError {
	return nil
}

// var _ Component = new(IterationComponent)

func (i *IterationComponent) Clear() {
	iterationComponentPool.Put(i)
}
