package core

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	integer = "integer"
	number  = "number"
	boolean = "boolean"
	array   = "array"
	object  = "object"
	str     = "string"
)

// ProcessNodeOutput 处理节点的所有输出
func ProcessNodeOutput(data map[string]any, outputs NodeDataOutputs) (map[string]any, error) {
	// 将数据转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("数据序列化失败: %v", err)
	}
	result := gjson.ParseBytes(jsonData)

	// 创建结果对象
	output := make(map[string]any)

	// 处理每个属性
	for name, prop := range outputs.Properties {
		// 获取输入值
		value := result.Get(name)

		// 如果值不存在且是必需字段
		if !value.Exists() {
			if isStringInSlice(name, outputs.Required) || prop.IsPropertyRequired {
				if prop.Default != nil {
					output[name] = prop.Default
					continue
				}
				return nil, fmt.Errorf("缺少必需字段: %s", name)
			}
			continue
		}

		// 根据类型处理值
		processedValue, err := processValue(value, prop)
		if err != nil {
			return nil, fmt.Errorf("处理字段 %s 失败: %w", name, err)
		}
		output[name] = processedValue
	}

	logx.Debugf("[输出处理结束] output:%+v", output)
	return output, nil
}

// processValue 根据属性定义处理值
func processValue(value gjson.Result, prop Properties) (any, error) {
	switch prop.Type {
	case str:
		if value.Type != gjson.String {
			return nil, fmt.Errorf("期望字符串类型，实际为 %s", value.Type.String())
		}
		return value.String(), nil

	case integer, number:
		if value.Type != gjson.Number {
			return nil, fmt.Errorf("期望数字类型，实际为 %s", value.Type.String())
		}
		if prop.Type == integer && value.Float() != float64(value.Int()) {
			return nil, fmt.Errorf("期望整数类型，实际为浮点数 %f", value.Float())
		}
		return value.Float(), nil

	case boolean:
		if value.Type != gjson.True && value.Type != gjson.False {
			return nil, fmt.Errorf("期望布尔类型，实际为 %s", value.Type.String())
		}
		return value.Bool(), nil

	case array:
		if !value.IsArray() {
			return nil, fmt.Errorf("期望数组类型，实际为 %s", value.Type.String())
		}

		// 处理数组元素
		arr := value.Array()
		result := make([]any, len(arr))
		for i, item := range arr {
			if prop.Item.Type == "" {
				// 如果没有指定元素类型，直接使用原始值
				result[i] = item.Value()
				continue
			}

			// 创建临时属性定义来处理数组元素
			itemProp := Properties{
				Type: prop.Item.Type,
				Item: prop.Item,
			}
			processedItem, err := processValue(item, itemProp)
			if err != nil {
				return nil, fmt.Errorf("处理数组元素 %d 失败: %w", i, err)
			}
			result[i] = processedItem
		}
		return result, nil

	case object:
		logx.Debugf("[输出处理] 对象字段")
		if !value.IsObject() {
			return nil, fmt.Errorf("期望对象类型，实际为 %s", value.Type.String())
		}

		// 处理对象属性
		result := make(map[string]any)
		for name, fieldProp := range prop.Properties {
			fieldValue := value.Get(name)
			if !fieldValue.Exists() {
				if isStringInSlice(name, prop.Item.Required) || fieldProp.IsPropertyRequired {
					if fieldProp.Default != nil {
						result[name] = fieldProp.Default
						continue
					}
					return nil, fmt.Errorf("对象缺少必需字段: %s", name)
				}
				continue
			}

			processedField, err := processValue(fieldValue, fieldProp)
			logx.Debugf("[输出处理] 对象字段 %s 处理结果: %+v", name, processedField)
			if err != nil {
				return nil, fmt.Errorf("处理对象字段 %s 失败: %w", name, err)
			}
			result[name] = processedField
		}
		return result, nil

	default:
		return nil, fmt.Errorf("不支持的类型: %s", prop.Type)
	}
}

// Helper
func isStringInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
