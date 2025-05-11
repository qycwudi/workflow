package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
)

// ValidateOutput 验证输出数据是否符合Schema定义
func ValidateOutput(data any, output Output) error {
	switch output.Type[0] {
	case "string":
		if _, ok := data.(string); !ok {
			logx.Errorf("[输出处理] 类型验证失败 [字段:%s] [期望:字符串] [实际:%s]", output.Name, reflect.TypeOf(data).String())
			return errors.New("输出 " + output.Name + " 必须是字符串类型")
		}
	case "integer":
		switch v := data.(type) {
		case int:
			if float64(v) != float64(v) {
				logx.Errorf("[输出处理] 整数验证失败 [字段:%s] [值:%v] [错误:不能有小数部分]", output.Name, v)
				return errors.New("输出 " + output.Name + " int必须是整数类型，不能有小数部分")
			}
		case int32:
			if float64(v) != float64(v) {
				logx.Errorf("[输出处理] 整数验证失败 [字段:%s] [值:%v] [错误:不能有小数部分]", output.Name, v)
				return errors.New("输出 " + output.Name + " int32必须是整数类型，不能有小数部分")
			}
		case int64:
			if float64(v) != float64(v) {
				logx.Errorf("[输出处理] 整数验证失败 [字段:%s] [值:%v] [错误:不能有小数部分]", output.Name, v)
				return errors.New("输出 " + output.Name + " int64必须是整数类型，不能有小数部分")
			}
		case float64:
			if v != float64(int(v)) {
				logx.Errorf("[输出处理] 整数验证失败 [字段:%s] [值:%v] [错误:不能有小数部分]", output.Name, v)
				return errors.New("输出 " + output.Name + " float64必须是整数类型，不能有小数部分")
			}
		case json.Number:
			_, err := strconv.ParseInt(string(v), 10, 64)
			if err != nil {
				logx.Errorf("[输出处理] 整数验证失败 [字段:%s] [值:%v] [错误:%v]", output.Name, v, err)
				return errors.New("输出 " + output.Name + " json.Number必须是整数类型")
			}
		default:
			logx.Errorf("[输出处理] 类型验证失败 [字段:%s] [期望:整数] [实际:%s]", output.Name, reflect.TypeOf(v).String())
			return errors.New("输出 " + output.Name + " ,类型: " + reflect.TypeOf(v).String() + " 必须是整数类型")
		}
	case "float":
		switch v := data.(type) {
		case float64:
			// float64 类型直接通过验证，允许有小数部分
			return nil
		case json.Number:
			_, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				logx.Errorf("[输出处理] 浮点数验证失败 [字段:%s] [值:%v] [错误:%v]", output.Name, v, err)
				return errors.New("输出 " + output.Name + " json.Number必须是浮点数类型")
			}
		default:
			logx.Errorf("[输出处理] 类型验证失败 [字段:%s] [期望:浮点数] [实际:%s]", output.Name, reflect.TypeOf(v).String())
			return errors.New("输出 " + output.Name + " ,类型: " + reflect.TypeOf(v).String() + " 必须是浮点数类型")
		}
	case "boolean":
		if _, ok := data.(bool); !ok {
			logx.Errorf("[输出处理] 类型验证失败 [字段:%s] [期望:布尔值] [实际:%s]", output.Name, reflect.TypeOf(data).String())
			return errors.New("输出 " + output.Name + " 必须是布尔类型")
		}
	case "array":
		arr, ok := data.([]any)
		if !ok {
			logx.Errorf("[输出处理] 类型验证失败 [字段:%s] [期望:数组] [实际:%s]", output.Name, reflect.TypeOf(data).String())
			return errors.New("输出 " + output.Name + " 必须是数组类型")
		}
		// 如果有更复杂的数组元素验证，可以在这里添加
		if output.Type[1] != "" {
			if len(arr) > 0 {
				err := ValidateOutput(arr[0], Output{Type: []string{output.Type[1]}})
				if err != nil {
					return errors.New("数组 " + output.Name + " 的元素 " + fmt.Sprintf("%+v", arr[0]) + " 验证失败: " + err.Error())
				}
			}
		}
	case "object":
		obj, ok := data.(map[string]any)
		if !ok {
			logx.Errorf("[输出处理] 类型验证失败 [字段:%s] [期望:对象] [实际:%+v]", output.Name, data)
			return errors.New("输出 " + output.Name + " 必须是对象类型,现在: " + fmt.Sprintf("%+v", data))
		}

		// 验证对象的每个字段
		for _, field := range output.Schema {
			fieldValue, exists := obj[field.Name]
			if !exists {
				logx.Errorf("[输出处理] 缺少必要字段 [对象:%s] [字段:%s]", output.Name, field.Name)
				return errors.New("对象 " + output.Name + " 缺少必要字段 " + field.Name)
			}

			if err := ValidateOutput(fieldValue, field); err != nil {
				logx.Errorf("[输出处理] 字段验证失败 [对象:%s] [字段:%s] [错误:%v]", output.Name, field.Name, err)
				return errors.New("对象 " + output.Name + " 的字段 " + field.Name + " 验证失败: " + err.Error())
			}
		}
	}
	return nil
}

// ParseOutput 解析单个输出
func ParseOutput(data any, output Output) (any, error) {
	if err := ValidateOutput(data, output); err != nil {
		logx.Errorf("[输出处理] 输出验证失败 [字段:%s] [错误:%v]", output.Name, err)
		return nil, err
	}

	switch output.Type[0] {
	case "object":
		obj := data.(map[string]any)
		result := make(map[string]any)
		if len(output.Schema) == 0 {
			return obj, nil
		}
		for _, field := range output.Schema {
			fieldValue, exists := obj[field.Name]
			if !exists {
				// 使用默认值
				result[field.Name] = field.DeftValue
				continue
			}

			parsedValue, err := ParseOutput(fieldValue, field)
			if err != nil {
				return nil, err
			}
			result[field.Name] = parsedValue
		}
		return result, nil
	case "array":
		// 简单数组直接返回
		return data, nil
	default:
		// 简单类型直接返回
		return data, nil
	}
}

// ProcessNodeOutput 处理节点的所有输出
func ProcessNodeOutput(data map[string]any, outputs []Output) (map[string]any, error) {
	var panicErr error
	result := make(map[string]any)

	// 处理 panics
	defer func() {
		if r := recover(); r != nil {
			panicErr = errors.New("ProcessNodeOutput panic: " + fmt.Sprintf("%+v", r))
			logx.Errorf("[输出处理] 处理过程发生panic [错误:%v]", r)
			// 清空结果，确保不会返回不完整的数据
			result = nil
		}
	}()

	if len(data) == 0 {
		return nil, nil
	}

	for _, output := range outputs {
		value, exists := data[output.Name]
		if !exists {
			// 如果输入中没有对应的数据，使用默认值或零值代替
			if output.Type[0] == "object" {
				r := make(map[string]any)
				err := sonic.Unmarshal([]byte(output.DeftValue.(string)), &r)
				if err != nil {
					logx.Errorf("[输出处理] 处理对象类输出失败 [字段:%s] [错误:%v]", output.Name, err)
					result[output.Name] = r // 保持原有行为，即使解析失败也返回解析后的结果
				} else {
					result[output.Name] = r
				}
			} else {
				switch output.Type[0] {
				case "integer":
					result[output.Name] = int64(0)
				case "float":
					result[output.Name] = float64(0)
				case "boolean":
					result[output.Name] = false
				case "array":
					result[output.Name] = make([]any, 0)
				case "object":
					result[output.Name] = make(map[string]any)
				case "string":
					result[output.Name] = ""
				default:
					result[output.Name] = output.DeftValue
				}
			}
			logx.Debugf("[输出处理] 使用默认值 [字段:%s] [值:%v]", output.Name, output.DeftValue)
			continue
		}

		parsedValue, err := ParseOutput(value, output)
		if err != nil {
			return nil, errors.New("处理输出 " + output.Name + " 时出错: " + err.Error())
		}

		result[output.Name] = parsedValue
		logx.Debugf("[输出处理] 输出解析成功 [字段:%s] [类型:%s]", output.Name, output.Type)
	}

	// 检查是否发生了 panic
	if panicErr != nil {
		return nil, panicErr
	}

	return result, nil
}
