package errors

import (
	"fmt"
	"time"
)

// ErrorType 错误类型枚举
type ErrorType string

const (
	TypeValidation    ErrorType = "validation"
	TypeExecution     ErrorType = "execution"
	TypeTimeout       ErrorType = "timeout"
	TypeNetwork       ErrorType = "network"
	TypeConfiguration ErrorType = "configuration"
	TypeCompilation   ErrorType = "compilation"
	TypeNotFound      ErrorType = "not_found"
	TypePermission    ErrorType = "permission"
)

// Severity 错误严重级别
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// ComponentError 组件错误结构
type ComponentError struct {
	Type        ErrorType         `json:"type"`
	Code        string            `json:"code"`
	Message     string            `json:"message"`
	Cause       error             `json:"-"`
	Component   string            `json:"component"`
	NodeID      string            `json:"node_id,omitempty"`
	WorkflowID  string            `json:"workflow_id,omitempty"`
	TraceID     string            `json:"trace_id,omitempty"`
	Severity    Severity          `json:"severity"`
	Timestamp   time.Time         `json:"timestamp"`
	Context     map[string]string `json:"context,omitempty"`
	Retryable   bool              `json:"retryable"`
	Suggestions []string          `json:"suggestions,omitempty"`
}

// Error 实现 error 接口
func (e *ComponentError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Component, e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Component, e.Type, e.Message)
}

// Unwrap 支持错误链
func (e *ComponentError) Unwrap() error {
	return e.Cause
}

// Is 实现错误类型判断
func (e *ComponentError) Is(target error) bool {
	if ce, ok := target.(*ComponentError); ok {
		return e.Type == ce.Type && e.Code == ce.Code
	}
	return false
}

// AddContext 添加上下文信息
func (e *ComponentError) AddContext(key, value string) *ComponentError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// AddSuggestion 添加解决建议
func (e *ComponentError) AddSuggestion(suggestion string) *ComponentError {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// WithCause 添加原因错误
func (e *ComponentError) WithCause(cause error) *ComponentError {
	e.Cause = cause
	return e
}