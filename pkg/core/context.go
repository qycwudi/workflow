package core

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/opentracing/opentracing-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// ExecutionContext 工作流执行上下文
type ExecutionContext struct {
	context.Context
	WorkspaceId string
	TraceId     string
	variables   map[string]interface{}
	State       *WorkflowState
	Route       map[string]struct{}
	mu          sync.RWMutex
	IsTrace     bool
	tracer      opentracing.Tracer
	startTime   time.Time
	metrics     map[string]float64

	TotalNodes int64
	Expiration time.Time // 添加过期时间
}

type WorkflowState struct {
	Status   ExecutionStatus
	Progress float64
	Result   map[string]*NodeResult
	Errors   []error
}

type NodeResult struct {
	Input    any
	Output   any
	Route    []string
	NodeID   string
	NodeName string
	Duration int64 // ms
	Error    string
	Type     string
}

// ExecutionStatus 执行状态
type ExecutionStatus string

const (
	StatusPending   ExecutionStatus = "pending"
	StatusRunning   ExecutionStatus = "running"
	StatusPaused    ExecutionStatus = "paused"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
	StatusCanceled  ExecutionStatus = "canceled"
)

var execContextPool = sync.Pool{
	New: func() interface{} {
		return &ExecutionContext{}
	},
}

// NewExecutionContext 创建新的执行上下文
func NewExecutionContext(ctx context.Context, workspaceId string, serialID string, totalNodes int64, params map[string]any) *ExecutionContext {

	// 如果ctx是ExecutionContext，则读取全部参数（用于loop组件）
	var allVariables map[string]any
	if ect, ok := ctx.(*ExecutionContext); ok {
		// 读取全部参数
		allVariables = ect.GetAllVariable()
		logx.Debugf("[Workflow] allVariables:%+v", allVariables)
	}

	execCtx := execContextPool.Get().(*ExecutionContext)
	// 如果serialID以trace-开头，则认为是追踪
	if strings.HasPrefix(serialID, "trace-") {
		execCtx.IsTrace = true
	}

	// execCtx.Context = ctx
	execCtx.TraceId = serialID
	execCtx.WorkspaceId = workspaceId
	execCtx.TotalNodes = totalNodes
	execCtx.State = NewWorkflowState()
	execCtx.Route = make(map[string]struct{})
	execCtx.startTime = time.Now()
	execCtx.Expiration = time.Now().Add(5 * time.Minute) // 设置过期时间
	execCtx.metrics = make(map[string]float64)

	execCtx.variables = make(map[string]any, totalNodes+1)
	execCtx.variables["_zero"] = params
	for key, value := range allVariables {
		execCtx.SetVariable(key, value)
	}
	return execCtx
}

// ReleaseExecutionContext 释放执行上下文
func ReleaseExecutionContext(ctx *ExecutionContext) {
	logx.Debugf("释放执行上下文: %s", ctx.TraceId)
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	// 清理状态
	ctx.State = nil
	ctx.variables = make(map[string]any)
	execContextPool.Put(ctx) // 将上下文放回池中
}

// ContextOption 上下文配置选项
type ContextOption func(*ExecutionContext)

// WithTracer 设置追踪器
func WithTracer(tracer opentracing.Tracer) ContextOption {
	return func(ctx *ExecutionContext) {
		// No-op, as the context is now a context.Context
	}
}

// GetVariable 获取变量值
func (ctx *ExecutionContext) GetVariable(key string) (any, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	value, ok := ctx.variables[key]
	return value, ok
}

func (ctx *ExecutionContext) GetAllVariable() map[string]any {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.variables
}

// SetVariable 设置变量值
func (ctx *ExecutionContext) SetVariable(key string, value any) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.variables[key] = value
}

// CheckRoute 检查路由
func (ctx *ExecutionContext) CheckRoute(key []string) bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	for _, route := range key {
		_, ok := ctx.Route[route]
		if ok {
			return true
		}
	}
	return false
}

// Cancel 取消执行
func (ctx *ExecutionContext) Cancel() {
	// No-op, as the context is now a context.Context
}

// Done 实现 context.Context 接口
func (c *ExecutionContext) Done() <-chan struct{} {
	// 直接返回底层 Context 的 Done channel
	return c.Context.Done()
}

// Deadline 实现 context.Context 接口
func (c *ExecutionContext) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

// Err 实现 context.Context 接口
func (c *ExecutionContext) Err() error {
	return c.Context.Err()
}

// Value 实现 context.Context 接口
func (c *ExecutionContext) Value(key interface{}) interface{} {
	if key == "execution_context" {
		return c
	}
	return c.Context.Value(key)
}

func (ctx *ExecutionContext) SetTracer(tracer opentracing.Tracer) {
	ctx.tracer = tracer
}

func (ctx *ExecutionContext) SetError(nodeID string, err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.State.Errors = append(ctx.State.Errors, err)
}

func (ctx *ExecutionContext) SetNodeResult(nodeID string, result *NodeResult) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.State.Result[nodeID] = result
	for _, route := range result.Route {
		ctx.Route[nodeID+"_"+route] = struct{}{}
	}
}

func NewWorkflowState() *WorkflowState {
	return &WorkflowState{
		Status: StatusPending,
		Errors: make([]error, 0),
		Result: make(map[string]*NodeResult),
	}
}

// GetNodeResult 获取指定节点的执行结果
func (ctx *ExecutionContext) GetNodeResult(nodeID string) (*NodeResult, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if ctx.State == nil || ctx.State.Result == nil {
		return nil, false
	}
	state, exists := ctx.State.Result[nodeID]
	if !exists || state == nil {
		return nil, false
	}
	return state, true
}

// MarshalResult 安全地序列化State.Result
func (ctx *ExecutionContext) MarshalResult() ([]byte, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return sonic.Marshal(ctx.State.Result)
}
