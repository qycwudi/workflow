package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	
	"workflow/pkg/core"
)

// Validator 配置验证器接口
type Validator interface {
	Validate(config any) []core.ValidationError
}

// Rule 验证规则接口
type Rule interface {
	Validate(field string, value any) *core.ValidationError
}

// FieldValidator 字段验证器
type FieldValidator struct {
	rules map[string][]Rule
}

// NewFieldValidator 创建字段验证器
func NewFieldValidator() *FieldValidator {
	return &FieldValidator{
		rules: make(map[string][]Rule),
	}
}

// AddRule 添加验证规则
func (v *FieldValidator) AddRule(field string, rule Rule) *FieldValidator {
	v.rules[field] = append(v.rules[field], rule)
	return v
}

// Validate 验证配置
func (v *FieldValidator) Validate(config any) []core.ValidationError {
	var validationErrors []core.ValidationError
	
	configValue := reflect.ValueOf(config)
	if configValue.Kind() == reflect.Ptr {
		configValue = configValue.Elem()
	}
	
	if configValue.Kind() != reflect.Struct {
		return []core.ValidationError{{
			Field:   "config",
			Message: "配置必须是结构体类型",
		}}
	}
	
	configType := configValue.Type()
	
	for i := 0; i < configValue.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := configValue.Field(i)
		
		// 获取JSON标签作为字段名
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		}
		
		// 验证字段
		if rules, exists := v.rules[fieldName]; exists {
			for _, rule := range rules {
				if err := rule.Validate(fieldName, fieldValue.Interface()); err != nil {
					validationErrors = append(validationErrors, *err)
				}
			}
		}
	}
	
	return validationErrors
}

// 预定义的验证规则

// RequiredRule 必填规则
type RequiredRule struct{}

func (r *RequiredRule) Validate(field string, value any) *core.ValidationError {
	if value == nil {
		return &core.ValidationError{
			Field:   field,
			Message: "字段不能为空",
		}
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if v.String() == "" {
			return &core.ValidationError{
				Field:   field,
				Message: "字段不能为空",
			}
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() == 0 {
			return &core.ValidationError{
				Field:   field,
				Message: "字段不能为空",
			}
		}
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return &core.ValidationError{
				Field:   field,
				Message: "字段不能为空",
			}
		}
	}
	
	return nil
}

// URLRule URL格式规则
type URLRule struct{}

func (r *URLRule) Validate(field string, value any) *core.ValidationError {
	str, ok := value.(string)
	if !ok {
		return &core.ValidationError{
			Field:   field,
			Message: "字段必须是字符串类型",
		}
	}
	
	if str == "" {
		return nil // 空值由RequiredRule处理
	}
	
	// 简单的URL格式验证
	urlPattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(urlPattern, str)
	if !matched {
		return &core.ValidationError{
			Field:   field,
			Message: "字段必须是有效的URL格式",
		}
	}
	
	return nil
}

// HTTPMethodRule HTTP方法规则
type HTTPMethodRule struct{}

func (r *HTTPMethodRule) Validate(field string, value any) *core.ValidationError {
	str, ok := value.(string)
	if !ok {
		return &core.ValidationError{
			Field:   field,
			Message: "字段必须是字符串类型",
		}
	}
	
	if str == "" {
		return nil
	}
	
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	method := strings.ToUpper(str)
	
	for _, validMethod := range validMethods {
		if method == validMethod {
			return nil
		}
	}
	
	return &core.ValidationError{
		Field:   field,
		Message: fmt.Sprintf("字段必须是有效的HTTP方法，支持: %s", strings.Join(validMethods, ", ")),
	}
}

// RangeRule 范围规则
type RangeRule struct {
	Min int64
	Max int64
}

func (r *RangeRule) Validate(field string, value any) *core.ValidationError {
	var num int64
	
	switch v := value.(type) {
	case int:
		num = int64(v)
	case int32:
		num = int64(v)
	case int64:
		num = v
	case float32:
		num = int64(v)
	case float64:
		num = int64(v)
	case string:
		var err error
		num, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return &core.ValidationError{
				Field:   field,
				Message: "字段必须是数字类型",
			}
		}
	default:
		return &core.ValidationError{
			Field:   field,
			Message: "字段必须是数字类型",
		}
	}
	
	if num < r.Min || num > r.Max {
		return &core.ValidationError{
			Field:   field,
			Message: fmt.Sprintf("字段值必须在 %d 到 %d 之间", r.Min, r.Max),
		}
	}
	
	return nil
}

// OneOfRule 选项规则
type OneOfRule struct {
	Options []string
}

func (r *OneOfRule) Validate(field string, value any) *core.ValidationError {
	str, ok := value.(string)
	if !ok {
		return &core.ValidationError{
			Field:   field,
			Message: "字段必须是字符串类型",
		}
	}
	
	if str == "" {
		return nil
	}
	
	for _, option := range r.Options {
		if str == option {
			return nil
		}
	}
	
	return &core.ValidationError{
		Field:   field,
		Message: fmt.Sprintf("字段值必须是以下选项之一: %s", strings.Join(r.Options, ", ")),
	}
}

// 便捷的构造函数

// Required 创建必填规则
func Required() Rule {
	return &RequiredRule{}
}

// URL 创建URL规则
func URL() Rule {
	return &URLRule{}
}

// HTTPMethod 创建HTTP方法规则
func HTTPMethod() Rule {
	return &HTTPMethodRule{}
}

// Range 创建范围规则
func Range(min, max int64) Rule {
	return &RangeRule{Min: min, Max: max}
}

// OneOf 创建选项规则
func OneOf(options ...string) Rule {
	return &OneOfRule{Options: options}
}