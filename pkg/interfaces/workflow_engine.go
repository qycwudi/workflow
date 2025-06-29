package interfaces

import (
	"context"
	"time"

	"workflow/pkg/core"
)

// WorkflowEngine 工作流引擎接口
type WorkflowEngine interface {
	// 编译和注册工作流
	Compile(ctx context.Context, id string, graph *core.Graph) error

	// 执行工作流
	ExecuteWorkflow(ctx context.Context, workflowID string, traceID string, params map[string]any, extra core.ContextExtra) (*core.ExecutionContextEnhanced, error)

	// 执行单个节点
	ExecuteSingleWorkflow(ctx context.Context, workflowID, serialID string, nodeId string, params map[string]any) (*core.NodeResult, error)

	// 暂停工作流
	PauseWorkflow(ctx context.Context, workflowID string, traceID string) error

	// 列出所有工作流
	ListWorkflows() []string

	// 移除工作流
	RemoveWorkflow(workflowID string) error

	// 获取工作流状态
	GetWorkflowStatus(workflowID string) (core.WorkflowStatus, error)

	// 关闭引擎
	Close() error
}

// Component 组件接口
type Component interface {
	// 验证组件配置
	Validate() []core.ValidationError

	// 分析输入参数
	AnalyzeInputs(ctx context.Context) (any, error)

	// 执行组件
	Execute(ctx context.Context, input any) (*core.Result, error)

	// 清理资源
	Clear()

	// 获取组件名称
	Name() string

	// 获取组件版本
	Version() string

	// 获取组件元数据
	Metadata() ComponentMetadata
}

// ComponentMetadata 组件元数据
type ComponentMetadata struct {
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author"`
	Category     string                 `json:"category"`
	Tags         []string               `json:"tags"`
	InputSchema  map[string]interface{} `json:"input_schema,omitempty"`
	OutputSchema map[string]interface{} `json:"output_schema,omitempty"`
	ConfigSchema map[string]interface{} `json:"config_schema,omitempty"`
}

// ExecutionContext 执行上下文接口
type ExecutionContext interface {
	context.Context

	// 变量管理
	GetVariable(key string) (map[string]any, bool)
	SetVariable(key string, value map[string]any)

	// 路由管理
	SetRoute(nodeID string, route []string)
	CheckRoute(key []string) bool

	// 生命周期管理
	Cancel()
	AddCleanupFunc(cleanup func())

	// 状态查询
	IsExpired() bool
	GetExecutionDuration() time.Duration

	// 元信息
	GetWorkspaceID() string
	GetTraceID() string
	GetExtra() core.ContextExtra
}

// TraceWriter trace写入器接口
type TraceWriter interface {
	// 创建trace记录
	CreateTrace(ctx context.Context, record TraceRecord) (int64, error)

	// 更新trace记录
	UpdateTrace(ctx context.Context, record TraceRecord) error

	// 查询trace记录
	QueryTraces(ctx context.Context, query TraceQuery) ([]TraceRecord, error)

	// 删除trace记录
	DeleteTrace(ctx context.Context, traceID string) error
}

// TraceRecord trace记录
type TraceRecord struct {
	ID          int64       `json:"id"`
	WorkspaceID string      `json:"workspace_id"`
	TraceID     string      `json:"trace_id"`
	NodeID      string      `json:"node_id"`
	NodeName    string      `json:"node_name"`
	NodeType    string      `json:"node_type"`
	Input       interface{} `json:"input"`
	Output      interface{} `json:"output"`
	Status      string      `json:"status"`
	StartTime   int64       `json:"start_time"`
	EndTime     int64       `json:"end_time"`
	Duration    int64       `json:"duration"`
	ErrorMsg    string      `json:"error_msg,omitempty"`
}

// TraceQuery trace查询条件
type TraceQuery struct {
	WorkspaceID string `json:"workspace_id,omitempty"`
	TraceID     string `json:"trace_id,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
	NodeType    string `json:"node_type,omitempty"`
	Status      string `json:"status,omitempty"`
	StartTime   int64  `json:"start_time,omitempty"`
	EndTime     int64  `json:"end_time,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Offset      int    `json:"offset,omitempty"`
}

// ConfigValidator 配置验证器接口
type ConfigValidator interface {
	// 验证配置
	Validate(config interface{}) []core.ValidationError

	// 添加验证规则
	AddRule(field string, rule ValidationRule) ConfigValidator

	// 获取验证规则
	GetRules() map[string][]ValidationRule
}

// ValidationRule 验证规则接口
type ValidationRule interface {
	// 验证字段值
	Validate(field string, value interface{}) *core.ValidationError

	// 获取规则名称
	Name() string

	// 获取规则描述
	Description() string
}

// PoolManager 协程池管理器接口
type PoolManager interface {
	// 提交任务
	Submit(task func()) error

	// 获取池统计信息
	GetStats() PoolStats

	// 调整池大小
	Resize(size int) error

	// 关闭池
	Close() error
}

// PoolStats 协程池统计信息
type PoolStats struct {
	Running  int `json:"running"`
	Free     int `json:"free"`
	Capacity int `json:"capacity"`
	Waiting  int `json:"waiting"`
}

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	// 记录执行时间
	RecordExecutionTime(workflowID, nodeID string, duration int64)

	// 记录执行次数
	RecordExecutionCount(workflowID, nodeID string, status string)

	// 记录错误
	RecordError(workflowID, nodeID string, errorType string)

	// 获取指标
	GetMetrics() map[string]interface{}

	// 重置指标
	ResetMetrics()
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// EventBus 事件总线接口
type EventBus interface {
	// 发布事件
	Publish(event Event) error

	// 订阅事件
	Subscribe(eventType string, handler EventHandler) error

	// 取消订阅
	Unsubscribe(eventType string, handler EventHandler) error

	// 关闭事件总线
	Close() error
}

// Event 事件接口
type Event interface {
	// 获取事件类型
	Type() string

	// 获取事件数据
	Data() interface{}

	// 获取事件时间戳
	Timestamp() int64

	// 获取事件ID
	ID() string
}

// EventHandler 事件处理器接口
type EventHandler interface {
	// 处理事件
	Handle(event Event) error

	// 获取处理器名称
	Name() string
}
