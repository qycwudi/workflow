package components

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/core"
)

// StartComponent 启动组件
type ConditionComponent struct {
	config ConditionConfig
}

type ConditionConfig struct {
	Conditions []core.Condition `json:"conditions"`
}

const (
	OperatorEquals      = "eq"           // 等于
	OperatorNotEquals   = "neq"          // 不等于
	OperatorGt          = "gt"           // 大于
	OperatorGte         = "gte"          // 大于等于
	OperatorLt          = "lt"           // 小于
	OperatorLte         = "lte"          // 小于等于
	OperatorIn          = "in"           // 包含 左值在右值数组中
	OperatorNotIn       = "nin"          // 不包含 左值不在右值数组中
	OperatorContains    = "contains"     // 字符串包含 左值包含右值
	OperatorNotContains = "ncontains"    // 字符串不包含 左值不包含右值
	OperatorIsEmpty     = "is_empty"     // 为空
	OperatorIsNotEmpty  = "is_not_empty" // 不为空
	OperatorIsTrue      = "is_true"      // 为真
	OperatorIsFalse     = "is_false"     // 为假
)

var branchComponentPool = sync.Pool{
	New: func() interface{} {
		return &ConditionComponent{}
	},
}

func (c *ConditionComponent) Name() string {
	return Condition
}

func NewConditionComponent(config json.RawMessage) (*ConditionComponent, error) {
	var branchConfig ConditionConfig
	var conditions []core.Condition
	if err := sonic.Unmarshal(config, &conditions); err != nil {
		return nil, errors.New("解析分支组件配置失败: " + err.Error())
	}
	// 使用pool
	component := branchComponentPool.Get().(*ConditionComponent)
	branchConfig.Conditions = conditions
	component.config = branchConfig
	return component, nil
}

func (c *ConditionComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	inputMap := input.(map[string]any)
	var route = []string{}
	for _, condition := range c.config.Conditions {
		left, ok := inputMap[condition.Key+"_left"]
		if !ok {
			logx.Errorf("左值不存在")
			continue
		}
		right, ok := inputMap[condition.Key+"_right"]
		if !ok {
			right = nil
		}
		logx.Debugf("branch execute expect: %s, condition: %s, left:%v, right:%v\n", condition.Key, condition.Value.Operator, left, right)
		switch condition.Value.Operator {
		case OperatorEquals:
			if left == right {
				route = append(route, condition.Key)
			}
		case OperatorNotEquals:
			if left != right {
				route = append(route, condition.Key)
			}
		case OperatorGt:
			// 比较数值
			left, right, err := compareValues(left, right)
			if err != nil {
				logx.Errorf("比较数值失败: %v", err)
				continue
			}
			if left <= right {
				continue
			}
			route = append(route, condition.Key)
		case OperatorGte:
			// 比较数值
			left, right, err := compareValues(left, right)
			if err != nil {
				logx.Errorf("比较数值失败: %v", err)
				continue
			}
			if left < right {
				continue
			}
			route = append(route, condition.Key)
		case OperatorLt:
			// 比较数值
			left, right, err := compareValues(left, right)
			if err != nil {
				logx.Errorf("比较数值失败: %v", err)
				continue
			}
			if left >= right {
				continue
			}
			route = append(route, condition.Key)
		case OperatorLte:
			// 比较数值
			left, right, err := compareValues(left, right)
			if err != nil {
				logx.Errorf("比较数值失败: %v", err)
				continue
			}
			if left > right {
				continue
			}
			route = append(route, condition.Key)
		case OperatorIn:
			// 检查左值是否在右值数组中
			if right == nil {
				logx.Errorf("IN 右值为空")
				continue
			}
			rightArray, ok := right.([]any)
			if !ok {
				logx.Errorf("IN 值类型不匹配,不是数组")
				continue
			}

			// 如果左值是单个值，直接检查是否在右值数组中
			if !slices.Contains(rightArray, left) {
				continue
			}

			route = append(route, condition.Key)
		case OperatorNotIn:
			// 检查左值是否不在右值数组中
			if right == nil {
				logx.Errorf("NOT IN 右值为空")
				continue
			}
			rightArray, ok := right.([]any)
			if !ok {
				logx.Errorf("NOT IN 右值类型不匹配,不是数组")
				continue
			}

			// 如果左值是单个值，直接检查是否不在右值数组中
			if slices.Contains(rightArray, left) {
				continue
			}
			route = append(route, condition.Key)
		case OperatorContains:
			leftStr, ok := left.(string)
			if !ok {
				logx.Errorf("左值类型不匹配,不是字符串")
				continue
			}
			rightStr, ok := right.(string)
			if !ok {
				logx.Errorf("右值类型不匹配,不是字符串")
				continue
			}
			if !strings.Contains(leftStr, rightStr) {
				continue
			}
			route = append(route, condition.Key)
		case OperatorNotContains:
			leftStr, ok := left.(string)
			if !ok {
				logx.Errorf("左值类型不匹配,不是字符串")
				continue
			}
			rightStr, ok := right.(string)
			if !ok {
				logx.Errorf("右值类型不匹配,不是字符串")
				continue
			}
			if strings.Contains(leftStr, rightStr) {
				continue
			}
			route = append(route, condition.Key)
		case OperatorIsEmpty:
			if left != nil {
				continue
			}
			route = append(route, condition.Key)
		case OperatorIsNotEmpty:
			if left == nil {
				continue
			}
			route = append(route, condition.Key)
		case OperatorIsTrue:
			leftBool, ok := left.(bool)
			if !ok {
				logx.Errorf("左值类型不匹配,不是布尔值")
				continue
			}
			if !leftBool {
				continue
			}
			route = append(route, condition.Key)
		case OperatorIsFalse:
			leftBool, ok := left.(bool)
			if !ok {
				logx.Errorf("左值类型不匹配,不是布尔值")
				continue
			}
			if leftBool {
				continue
			}
			route = append(route, condition.Key)
		}
		logx.Debugf("branch execute condition: %s, route:%v\n", condition.Value.Operator, route)
	}
	if len(route) == 0 {
		logx.Debugf("没有满足条件的路由,执行else路由")
		route = append(route, Else)
	}
	outPut := map[string]interface{}{
		"routes": route,
	}
	return &core.Result{
		Route:  route,
		Output: outPut,
	}, nil
}

func (c *ConditionComponent) Validate() []core.ValidationError {
	if len(c.config.Conditions) == 0 {
		return []core.ValidationError{
			{
				Field:   "conditions",
				Message: "条件不能为空",
			},
		}
	}
	for _, condition := range c.config.Conditions {
		if condition.Value.Operator == "" {
			return []core.ValidationError{
				{
					Field:   "route",
					Message: "操作符不能为空",
				},
			}
		}
	}
	return nil
}
func (c *ConditionComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	execCtx := ctx.(*core.ExecutionContext)
	result := make(map[string]any)
	for i, condition := range c.config.Conditions {
		// 构造输入值
		leftKey := condition.Key + "_left"
		leftInputValue := core.NodeDataInputsValues{
			Type:    "ref",
			Content: condition.Value.Left.Content,
		}
		leftInputMap := map[string]core.NodeDataInputsValues{
			leftKey: leftInputValue,
		}

		leftInput := core.NodeDataInputs{
			Properties: map[string]core.Properties{
				leftKey: {
					Type: condition.LeftType,
				},
			},
			Required: []string{leftKey},
		}
		left, err := core.ParseNodeInputs(execCtx, leftInputMap, leftInput)
		if err != nil {
			return nil, err
		}
		result[leftKey] = left[leftKey]

		// 构造右值输入值
		var rightInputValue core.NodeDataInputsValues
		rightKey := condition.Key + "_right"

		if condition.RightType != "undefined" {
			rightInputValue = core.NodeDataInputsValues{
				Type:    "ref",
				Content: condition.Value.Right.Content,
			}

			rightInputMap := map[string]core.NodeDataInputsValues{
				rightKey: rightInputValue,
			}
			rightInput := core.NodeDataInputs{
				Properties: map[string]core.Properties{
					rightKey: {
						Type: condition.RightType,
					},
				},
				Required: []string{rightKey},
			}
			right, err := core.ParseNodeInputs(execCtx, rightInputMap, rightInput)
			if err != nil {
				return nil, err
			}
			result[rightKey] = right[rightKey]
		} else {
			if condition.Value.Right.Type == core.CONTENT {
				// 要和左边的值类型保持一致 leftType
				convertedValue, err := core.ConvertValue(condition.Value.Right.Content, condition.LeftType)
				if err != nil {
					return nil, err
				}
				result[rightKey] = convertedValue
			}
		}
		logx.Debugf("branch component analyze inputs values: index:%d, %+v\n", i, result)
	}
	return result, nil
}

// compareValues 比较两个值
func compareValues(left, right any) (float64, float64, error) {
	// 转换为float64进行比较
	var leftFloat, rightFloat float64

	switch v := left.(type) {
	case int:
		leftFloat = float64(v)
	case int64:
		leftFloat = float64(v)
	case float32:
		leftFloat = float64(v)
	case float64:
		leftFloat = v
	default:
		return 0, 0, errors.New("不支持的左值类型: " + reflect.TypeOf(left).String())
	}

	switch v := right.(type) {
	case int:
		rightFloat = float64(v)
	case int64:
		rightFloat = float64(v)
	case float32:
		rightFloat = float64(v)
	case float64:
		rightFloat = v
	default:
		return 0, 0, errors.New("不支持的右值类型: " + reflect.TypeOf(right).String())
	}

	logx.Debugf("compareValues: left:%v, right:%v\n", leftFloat, rightFloat)
	return leftFloat, rightFloat, nil
}

func (c *ConditionComponent) Clear() {
	branchComponentPool.Put(c)
}
