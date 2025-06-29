package constants

import "time"

// Component Types - 组件类型常量
const (
	ComponentStart     = "start"
	ComponentEnd       = "end"
	ComponentHTTP      = "http"
	ComponentCode      = "code"
	ComponentCondition = "condition"
	ComponentModel     = "llm"
	ComponentLoop      = "loop"
	ComponentStartItem = "start-item"
	ComponentEndItem   = "end-item"
	ComponentDatabase  = "sql"
)

// Execution Routes - 执行路由常量
const (
	RouteSuccess = "success"
	RouteFailed  = "failed"
	RouteSkip    = "skip"
	RouteTrue    = "true"
	RouteFalse   = "false"
	RouteElse    = "else"
)

// Context Keys - 上下文键常量
const (
	GenesisParameters = "_zero_"
	EndParameters     = "_end_"
	TracePrefix       = "trace-"
)

// HTTP Component Constants
const (
	HTTPHeaderKey = "headers"
	HTTPBodyKey   = "body"
	HTTPURLKey    = "url"
	HTTPIsJSONKey = "isJSON"
	
	HTTPBodyTypeJSON     = "json"
	HTTPBodyTypeForm     = "form-data"
	HTTPBodyTypeNone     = "none"
)

// Loop Component Constants
const (
	LoopItemKey = "items"
)

// JavaScript Engine Constants
const (
	JSLanguage     = "javascript"
	JSType         = "Js"
	GolangLanguage = "golang"
	PythonLanguage = "python"
	
	ErrorHandlingModeAbort = "abort"
	ErrorHandlingModeRetry = "retry"
)

// Default Configuration Values
const (
	DefaultInitialPoolSize        = 100
	DefaultMaxPoolSize            = 10000
	DefaultMaxConcurrentWorkflows = 1000
	DefaultWorkflowTimeout        = 2 * time.Minute
	DefaultHTTPTimeout            = 30 * time.Second
	DefaultJSTimeout              = 30 // seconds
	DefaultHTTPRetries            = 3
)

// Error Codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeExecution      = "EXECUTION_ERROR"
	ErrCodeTimeout        = "TIMEOUT_ERROR"
	ErrCodeNetwork        = "NETWORK_ERROR"
	ErrCodeConfiguration  = "CONFIG_ERROR"
	ErrCodeCompilation    = "COMPILATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND_ERROR"
	ErrCodePermission     = "PERMISSION_ERROR"
)