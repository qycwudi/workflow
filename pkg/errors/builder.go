package errors

import (
	"fmt"
	"time"
	
	"workflow/pkg/constants"
)

// ErrorBuilder 错误构建器
type ErrorBuilder struct {
	err *ComponentError
}

// NewError 创建新的错误构建器
func NewError(errorType ErrorType, component string) *ErrorBuilder {
	return &ErrorBuilder{
		err: &ComponentError{
			Type:      errorType,
			Component: component,
			Timestamp: time.Now(),
			Retryable: false,
			Context:   make(map[string]string),
		},
	}
}

// Code 设置错误码
func (b *ErrorBuilder) Code(code string) *ErrorBuilder {
	b.err.Code = code
	return b
}

// Message 设置错误消息
func (b *ErrorBuilder) Message(message string) *ErrorBuilder {
	b.err.Message = message
	return b
}

// Messagef 设置格式化错误消息
func (b *ErrorBuilder) Messagef(format string, args ...interface{}) *ErrorBuilder {
	b.err.Message = fmt.Sprintf(format, args...)
	return b
}

// Cause 设置原因错误
func (b *ErrorBuilder) Cause(cause error) *ErrorBuilder {
	b.err.Cause = cause
	return b
}

// NodeID 设置节点ID
func (b *ErrorBuilder) NodeID(nodeID string) *ErrorBuilder {
	b.err.NodeID = nodeID
	return b
}

// WorkflowID 设置工作流ID
func (b *ErrorBuilder) WorkflowID(workflowID string) *ErrorBuilder {
	b.err.WorkflowID = workflowID
	return b
}

// TraceID 设置跟踪ID
func (b *ErrorBuilder) TraceID(traceID string) *ErrorBuilder {
	b.err.TraceID = traceID
	return b
}

// Severity 设置严重级别
func (b *ErrorBuilder) Severity(severity Severity) *ErrorBuilder {
	b.err.Severity = severity
	return b
}

// Context 添加上下文
func (b *ErrorBuilder) Context(key, value string) *ErrorBuilder {
	b.err.Context[key] = value
	return b
}

// Retryable 设置是否可重试
func (b *ErrorBuilder) Retryable(retryable bool) *ErrorBuilder {
	b.err.Retryable = retryable
	return b
}

// Suggestion 添加解决建议
func (b *ErrorBuilder) Suggestion(suggestion string) *ErrorBuilder {
	b.err.Suggestions = append(b.err.Suggestions, suggestion)
	return b
}

// Build 构建错误
func (b *ErrorBuilder) Build() *ComponentError {
	// 设置默认错误码
	if b.err.Code == "" {
		switch b.err.Type {
		case TypeValidation:
			b.err.Code = constants.ErrCodeValidation
		case TypeExecution:
			b.err.Code = constants.ErrCodeExecution
		case TypeTimeout:
			b.err.Code = constants.ErrCodeTimeout
		case TypeNetwork:
			b.err.Code = constants.ErrCodeNetwork
		case TypeConfiguration:
			b.err.Code = constants.ErrCodeConfiguration
		case TypeCompilation:
			b.err.Code = constants.ErrCodeCompilation
		case TypeNotFound:
			b.err.Code = constants.ErrCodeNotFound
		case TypePermission:
			b.err.Code = constants.ErrCodePermission
		}
	}
	
	// 设置默认严重级别
	if b.err.Severity == "" {
		switch b.err.Type {
		case TypeValidation, TypeConfiguration:
			b.err.Severity = SeverityMedium
		case TypeTimeout, TypeNetwork:
			b.err.Severity = SeverityHigh
		case TypeExecution, TypeCompilation:
			b.err.Severity = SeverityHigh
		default:
			b.err.Severity = SeverityMedium
		}
	}
	
	return b.err
}

// 预定义的错误构建器工厂函数

// ValidationError 创建验证错误
func ValidationError(component, message string) *ErrorBuilder {
	return NewError(TypeValidation, component).
		Message(message).
		Severity(SeverityMedium).
		Suggestion("请检查组件配置参数")
}

// ExecutionError 创建执行错误
func ExecutionError(component, message string) *ErrorBuilder {
	return NewError(TypeExecution, component).
		Message(message).
		Severity(SeverityHigh).
		Retryable(true)
}

// TimeoutError 创建超时错误
func TimeoutError(component, message string) *ErrorBuilder {
	return NewError(TypeTimeout, component).
		Message(message).
		Severity(SeverityHigh).
		Retryable(true).
		Suggestion("请检查网络连接或增加超时时间")
}

// NetworkError 创建网络错误
func NetworkError(component, message string) *ErrorBuilder {
	return NewError(TypeNetwork, component).
		Message(message).
		Severity(SeverityHigh).
		Retryable(true).
		Suggestion("请检查网络连接")
}

// ConfigurationError 创建配置错误
func ConfigurationError(component, message string) *ErrorBuilder {
	return NewError(TypeConfiguration, component).
		Message(message).
		Severity(SeverityMedium).
		Suggestion("请检查组件配置")
}

// CompilationError 创建编译错误
func CompilationError(component, message string) *ErrorBuilder {
	return NewError(TypeCompilation, component).
		Message(message).
		Severity(SeverityHigh).
		Suggestion("请检查代码语法")
}

// NotFoundError 创建未找到错误
func NotFoundError(component, message string) *ErrorBuilder {
	return NewError(TypeNotFound, component).
		Message(message).
		Severity(SeverityMedium)
}