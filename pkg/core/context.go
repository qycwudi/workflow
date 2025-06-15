package core

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// ExecutionContext 工作流执行上下文
type ExecutionContext struct {
	context.Context
	WorkspaceId     string
	TraceId         string
	parentVariables *KVStore
	variables       *KVStore
	Route           map[string]struct{}
	mu              sync.RWMutex
	IsTrace         bool
	startTime       time.Time
	Expiration      time.Time // 添加过期时间
	Extra           ContextExtra
	cancel          context.CancelFunc
}

type ContextExtra struct {
	Index             int64  // 索引
	IsSub             bool   // 是否是子流程
	ParentWorkspaceId string // 父流程workspaceId
	NodeNum           int64  // 节点数量
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

var GenesisParameters = "_zero_"
var EndParameters = "_end_"

var TracePrefix = "trace-"

var execContextPool = sync.Pool{
	New: func() interface{} {
		return &ExecutionContext{}
	},
}

// NewExecutionContext 创建新的执行上下文
func NewExecutionContext(ctx context.Context, workspaceId string, traceId string, params map[string]any, extra ContextExtra, timeout time.Duration, nodeNum int64) *ExecutionContext {

	execCtx := execContextPool.Get().(*ExecutionContext)
	// 如果serialID以trace-开头，则认为是追踪
	if strings.HasPrefix(traceId, TracePrefix) {
		execCtx.IsTrace = true
	}
	execCtx.Extra = extra
	withTimeout, cancel := context.WithTimeout(ctx, timeout)
	execCtx.Context = withTimeout
	execCtx.cancel = cancel
	execCtx.Expiration = time.Now().Add(timeout)
	execCtx.TraceId = traceId
	execCtx.WorkspaceId = workspaceId
	execCtx.Route = make(map[string]struct{})
	execCtx.startTime = time.Now()

	// 设置变量缓存大小 默认存储输入输出变量
	execCtx.variables = CreateCache()
	if extra.IsSub {
		if ect, ok := ctx.(*ExecutionContext); ok {
			execCtx.parentVariables = ect.variables
		} else {
			logx.Errorf("[Workflow] NewExecutionContext: ctx is not ExecutionContext,Failed to reuse variables")
		}
	}
	// 设置初始参数 workspaceId作为前缀防止迭代时变量覆盖
	execCtx.SetVariable(GenesisParameters+workspaceId, params)
	execCtx.Extra.NodeNum = nodeNum
	return execCtx
}

// ReleaseExecutionContext 释放执行上下文
func ReleaseExecutionContext(ctx *ExecutionContext) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	// 清理状态
	ctx.variables.Recover()
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
func (ctx *ExecutionContext) GetVariable(key string) (map[string]any, bool) {
	value, ok := ctx.variables.Get(key)
	if !ok {
		if ctx.parentVariables != nil {
			value, ok = ctx.parentVariables.Get(key)
		}
	}
	return value, ok
}

// SetVariable 设置变量值
func (ctx *ExecutionContext) SetVariable(key string, value map[string]any) {
	ctx.variables.Set(key, value)
}

// SetRoute 设置路由
func (ctx *ExecutionContext) SetRoute(nodeID string, route []string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	for _, r := range route {
		ctx.Route[nodeID+"_"+r] = struct{}{}
	}
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
	ctx.cancel()
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
