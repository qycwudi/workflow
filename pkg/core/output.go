package core

import (
	"encoding/json"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/rotisserie/eris"
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
		return nil, eris.New(fmt.Sprintf("ProcessNodeOutput Data Marshal failed: %v", err))
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
					switch prop.Type {
					case object:
						{
							// 强转 string
							if defaultJson, ok := prop.Default.(string); ok {
								var r map[string]any
								err := sonic.UnmarshalString(defaultJson, &r)
								if err != nil {
									return nil, eris.New(fmt.Sprintf("Failed to process object field %s: %s", name, err.Error()))
								}
								output[name] = r
							} else {
								// 默认值不是 string 类型，则直接赋空 map
								logx.Errorf("[Output processing] Object field %s default value is not string type: %v", name, prop.Default)
								output[name] = make(map[string]any)
							}
						}
					case str:
						{
							// 强转 string
							if defaultStr, ok := prop.Default.(string); ok {
								output[name] = defaultStr
							} else {
								output[name] = ""
							}
						}
					case number:
						{
							// 强转 number
							if defaultNum, ok := prop.Default.(float64); ok {
								output[name] = defaultNum
							} else {
								output[name] = 0
							}
						}
					case boolean:
						{
							// 强转 boolean
							if defaultBool, ok := prop.Default.(bool); ok {
								output[name] = defaultBool
							} else {
								output[name] = false
							}
						}
					case array:
						{
							// 强转 array
							if defaultArr, ok := prop.Default.([]any); ok {
								output[name] = defaultArr
							} else {
								output[name] = []any{}
							}
						}
					}
					continue
				}
				return nil, eris.New(fmt.Sprintf("Missing required field: %s", name))
			}
			continue
		}

		// 根据类型处理值
		processedValue, err := processValue(value, prop)
		if err != nil {
			logx.Errorf("[Output processing] Failed to process field %s: %s", name, err.Error())
			return nil, eris.New(fmt.Sprintf("Failed to process field %s: %s", name, err.Error()))
		}
		output[name] = processedValue
	}
	return output, nil
}

// processValue 根据属性定义处理值
func processValue(value gjson.Result, prop Properties) (any, error) {
	switch prop.Type {
	case str:
		if value.Type != gjson.String {
			return nil, eris.New(fmt.Sprintf("Expected string type, got %s", value.Type.String()))
		}
		return value.String(), nil

	case integer, number:
		if value.Type != gjson.Number {
			return nil, eris.New(fmt.Sprintf("Expected number type, got %s", value.Type.String()))
		}
		if prop.Type == integer && value.Float() != float64(value.Int()) {
			return nil, eris.New(fmt.Sprintf("Expected integer type, got float %f", value.Float()))
		}
		return value.Float(), nil

	case boolean:
		if value.Type != gjson.True && value.Type != gjson.False {
			return nil, eris.New(fmt.Sprintf("Expected boolean type, got %s", value.Type.String()))
		}
		return value.Bool(), nil

	case array:
		if !value.IsArray() {
			return nil, eris.New(fmt.Sprintf("Expected array type, got %s", value.Type.String()))
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
				return nil, eris.New(fmt.Sprintf("Failed to process array element %d: %s", i, err.Error()))
			}
			result[i] = processedItem
		}
		return result, nil

	case object:
		logx.Debugf("[Output processing] Object field")
		if !value.IsObject() {
			return nil, eris.New(fmt.Sprintf("Expected object type, got %s", value.Type.String()))
		}

		// 处理对象属性
		result := make(map[string]any)
		for name, fieldProp := range prop.Properties {
			fieldValue := value.Get(name)
			if !fieldValue.Exists() {
				if isStringInSlice(name, prop.Item.Required) || fieldProp.IsPropertyRequired {
					if fieldProp.Default != nil {
						// 强转 string
						if defaultJson, ok := fieldProp.Default.(string); ok {
							var r map[string]any
							err := sonic.UnmarshalString(defaultJson, &r)
							if err != nil {
								return nil, eris.New(fmt.Sprintf("Failed to process object field %s: %s", name, err.Error()))
							}
							result[name] = r
						} else {
							// 默认值不是 string 类型，则直接赋空 map
							logx.Errorf("[Output processing] Object field %s default value is not string type: %v", name, fieldProp.Default)
							result[name] = make(map[string]any)
						}

						continue
					}
					return nil, eris.New(fmt.Sprintf("Object missing required field: %s", name))
				}
				continue
			}

			processedField, err := processValue(fieldValue, fieldProp)
			logx.Debugf("[Output processing] Object field %s processing result: %+v", name, processedField)
			if err != nil {
				return nil, eris.New(fmt.Sprintf("Failed to process object field %s: %s", name, err.Error()))
			}
			result[name] = processedField
		}
		return result, nil

	default:
		return nil, eris.New(fmt.Sprintf("Unsupported type: %s", prop.Type))
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
