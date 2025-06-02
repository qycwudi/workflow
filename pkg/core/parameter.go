package core

// Result 组件执行结果
type Result struct {
	Output any
	Route  []string
}

// ValidationError 验证错误
type ValidationError struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Severity   string `json:"severity"`
	Suggestion string `json:"suggestion"`
}
