package core

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

// 从父节点输出数据中提取指定路径的数据
func extractDataByPath(data any, path string) (any, error) {
	if path == "" {
		return data, nil
	}

	// 分割路径
	parts := strings.Split(path, ".")
	current := data

	// 遍历路径的每个部分
	for _, part := range parts {
		// 处理 map 类型
		if mapData, ok := current.(map[string]any); ok {
			if value, exists := mapData[part]; exists {
				current = value
				continue
			}
			logx.Errorf("[输入处理] 字段不存在 [路径:%s] [字段:%s]", path, part)
			return nil, errors.New("路径 " + path + " 中的字段 " + part + " 不存在")
		}

		// 处理数组类型
		if arrayData, ok := current.([]any); ok {
			// 尝试将部分转换为索引
			index, err := strconv.Atoi(part)
			if err != nil {
				logx.Errorf("[输入处理] 数组索引解析失败 [路径:%s] [部分:%s] [错误:%v]", path, part, err)
				return nil, errors.New("路径 " + path + " 中的 " + part + " 不是有效的数组索引")
			}
			if index < 0 || index >= len(arrayData) {
				logx.Errorf("[输入处理] 数组索引越界 [路径:%s] [索引:%d] [长度:%d]", path, index, len(arrayData))
				return nil, errors.New("数组索引 " + strconv.Itoa(index) + " 超出范围")
			}
			current = arrayData[index]
			continue
		}

		logx.Errorf("[输入处理] 不支持的数据类型 [路径:%s] [字段:%s] [类型:%s]", path, part, reflect.TypeOf(current).String())
		return nil, errors.New("路径 " + path + " 中的 " + part + " 不支持的数据类型: " + reflect.TypeOf(current).String())
	}

	return current, nil
}

// 类型检查和转换
func validateAndConvertType(value any, expectedType []string) (any, error) {
	switch expectedType[0] {
	case "string":
		if str, ok := value.(string); ok {
			return str, nil
		}
		logx.Errorf("[输入处理] 类型不匹配 [期望:%s] [实际:%s]", expectedType[0], reflect.TypeOf(value).String())
		return nil, errors.New("类型不匹配: 期望 " + expectedType[0] + ", 实际是 " + reflect.TypeOf(value).String())

	case "integer":
		// 统一解析成 int64
		if num, ok := value.(float64); ok {
			return int64(num), nil
		}
		if num, ok := value.(int64); ok {
			return int64(num), nil
		}
		if num, ok := value.(int32); ok {
			return int64(num), nil
		}
		if num, ok := value.(int); ok {
			return int64(num), nil
		}
		if num, ok := value.(json.Number); ok {
			_, err := strconv.ParseInt(string(num), 10, 64)
			if err != nil {
				logx.Errorf("[输入处理] 整数类型转换失败 [值:%v] [错误:%v]", value, err)
				return nil, errors.New("类型不匹配: 期望 integer, 实际是 " + reflect.TypeOf(value).String())
			}
			return num, nil
		}
		return nil, errors.New("类型不匹配: 期望 integer, 实际是 " + reflect.TypeOf(value).String())

	case "float":
		if num, ok := value.(float64); ok {
			return num, nil
		}
		// 处理整数转浮点数的情况
		if num, ok := value.(int); ok {
			return float64(num), nil
		}
		if num, ok := value.(json.Number); ok {
			_, err := strconv.ParseFloat(string(num), 64)
			if err != nil {
				return nil, errors.New("类型不匹配: 期望 float, 实际是 " + reflect.TypeOf(value).String())
			}
			return num, nil
		}
		return nil, errors.New("类型不匹配: 期望 float, 实际是 " + reflect.TypeOf(value).String())

	case "boolean":
		if b, ok := value.(bool); ok {
			return b, nil
		}
		return nil, errors.New("类型不匹配: 期望 boolean, 实际是 " + reflect.TypeOf(value).String())

	case "array":
		if arr, ok := value.([]any); ok {
			// 解析expectedType[1]元素类型
			if len(expectedType) <= 1 {
				return arr, nil
			}
			elementType := expectedType[1]
			for _, item := range arr {
				_, err := validateAndConvertType(item, []string{elementType})
				if err != nil {
					logx.Errorf("[输入处理] 数组元素类型不匹配 [错误:%v]", err)
					return nil, errors.New("数组元素类型不匹配: " + err.Error())
				}
			}
			return arr, nil
		}
		return nil, errors.New("类型不匹配: 期望 array, 实际是 " + reflect.TypeOf(value).String())

	case "object":
		if obj, ok := value.(map[string]any); ok {
			return obj, nil
		}
		return nil, errors.New("类型不匹配: 期望 object, 实际是 " + reflect.TypeOf(value).String())

	default:
		return value, nil
	}
}

// ParseNodeInputs 解析节点输入
func ParseNodeInputs(inputs []Inputs, parentOutputs *ExecutionContext) (map[string]any, error) {
	result := make(map[string]any)

	for _, input := range inputs {
		logx.Debugf("[输入处理] 解析输入 [字段:%s] [类型:%s] [值类型:%s]", input.Name, input.Type, input.Value.Type)
		// 如果类型是fix，则直接写入值
		if input.Value.Type == "fix" {
			result[input.Name] = input.Value.Content.Value
			logx.Debugf("[输入处理] 固定输入解析成功 [字段:%s] [类型:%s]", input.Name, input.Type)
			continue
		}
		// 获取父节点ID和输出字段名
		parentID := input.Value.Content.BlockID
		outputName := input.Value.Content.Name

		// 获取父节点的输出数据
		parentOutput, exists := parentOutputs.GetVariable(parentID + ".output")
		if !exists {
			return nil, errors.New("找不到父节点 " + parentID + " 的输出数据")
		}
		// 从父节点输出中提取数据
		value, err := extractDataByPath(parentOutput, outputName)
		if err != nil {
			return nil, errors.New("解析输入 " + input.Name + " 时出错: " + err.Error())
		}

		// 验证和转换类型
		convertedValue, err := validateAndConvertType(value, input.Type)
		if err != nil {
			return nil, errors.New("输入 " + input.Name + " 的类型验证失败: " + err.Error())
		}

		// 将解析后的值添加到结果中
		result[input.Name] = convertedValue
		logx.Debugf("[输入处理] 输入解析成功 [字段:%s] [类型:%s] [值:%+v]", input.Name, input.Type, convertedValue)
	}

	return result, nil
}
