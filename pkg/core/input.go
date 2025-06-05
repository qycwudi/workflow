package core

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	REF     = "ref"
	CONTENT = "constant"
)

// ParseNodeInputs 解析节点输入
func ParseNodeInputs(parentOutputs *ExecutionContext, inputsValues map[string]NodeDataInputsValues, inputs NodeDataInputs) (map[string]any, error) {
	result := make(map[string]any)

	// 遍历所有输入值
	for name, inputValue := range inputsValues {
		// 检查是否是必需字段
		if isStringInSlice(name, inputs.Required) {
			if inputValue.Type == "" {
				return nil, fmt.Errorf("缺少必需字段: %s", name)
			}
		}

		// 获取字段定义
		prop, exists := inputs.Properties[name]
		if !exists {
			logx.Debugf("[输入处理] 未定义的输入字段: %s", name)
			continue
		}

		var value any
		var err error

		// 根据输入类型处理
		switch inputValue.Type {
		case CONTENT:
			// 常量类型直接使用值
			value = inputValue.Content
		case REF:
			// 引用类型需要从父节点获取
			content, ok := inputValue.Content.([]any)
			if !ok || len(content) < 2 {
				return nil, fmt.Errorf("引用格式错误: %v,ok:%v,len:%v", inputValue.Content, ok, len(content))
			}

			parentID := content[0].(string)

			fieldPath := ""
			for _, v := range content[1:] {
				fieldPath += "." + v.(string)
			}
			// 去掉最前面的.
			fieldPath = fieldPath[1:]
			// 获取父节点输出
			parentOutput, exists := parentOutputs.GetVariable(parentID + ".output")
			if !exists {
				return nil, fmt.Errorf("找不到父节点 %s 的输出数据", parentID)
			}
			// 将父节点输出转换为 JSON
			jsonData, err := json.Marshal(parentOutput)
			if err != nil {
				return nil, fmt.Errorf("父节点数据序列化失败: %v", err)
			}
			logx.Debugf("[输入处理] jsonData:%s,fieldPath:%s", string(jsonData), fieldPath)
			// 使用 gjson 获取指定路径的值
			value = gjson.GetBytes(jsonData, fieldPath).Value()
			if value == nil {
				return nil, fmt.Errorf("在父节点 %s 中找不到字段 %s", parentID, fieldPath)
			}
		default:
			return nil, fmt.Errorf("不支持的输入类型: %s", inputValue.Type)
		}

		// 类型转换和验证
		convertedValue, err := ConvertValue(value, prop.Type)
		if err != nil {
			return nil, fmt.Errorf("字段 %s 类型转换失败: %v", name, err)
		}

		result[name] = convertedValue
		logx.Debugf("[输入处理] 输入解析成功 [字段:%s] [类型:%s] [值:%+v]", name, prop.Type, convertedValue)
	}
	logx.Debugf("[输入处理结束] result:%+v", result)
	return result, nil
}

// convertValue 转换值到指定类型
func ConvertValue(value any, targetType string) (any, error) {
	// 将值转换为 JSON 字符串
	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("值序列化失败: %v", err)
	}
	result := gjson.ParseBytes(jsonData)

	switch targetType {
	case "string":
		if result.Type != gjson.String {
			return nil, fmt.Errorf("期望字符串类型，实际为 %s", result.Type.String())
		}
		return result.String(), nil

	case "number":
		if result.Type != gjson.Number {
			return nil, fmt.Errorf("期望数字类型，实际为 %s", result.Type.String())
		}
		return result.Float(), nil

	case "integer":
		if result.Type != gjson.Number {
			return nil, fmt.Errorf("期望整数类型，实际为 %s", result.Type.String())
		}
		if result.Float() != float64(result.Int()) {
			return nil, fmt.Errorf("期望整数类型，实际为浮点数 %f", result.Float())
		}
		return result.Int(), nil

	case "boolean":
		if result.Type != gjson.True && result.Type != gjson.False {
			return nil, fmt.Errorf("期望布尔类型，实际为 %s", result.Type.String())
		}
		return result.Bool(), nil

	case "array":
		if !result.IsArray() {
			return nil, fmt.Errorf("期望数组类型，实际为 %s", result.Type.String())
		}
		return result.Value(), nil

	case "object":
		if !result.IsObject() {
			return nil, fmt.Errorf("期望对象类型，实际为 %s", result.Type.String())
		}
		return result.Value(), nil

	default:
		return nil, fmt.Errorf("不支持的类型: %s", targetType)
	}
}
