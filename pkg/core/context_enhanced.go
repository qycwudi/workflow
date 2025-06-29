package core

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/constants"
)

// ExecutionContextEnhanced 增强的执行上下文
type ExecutionContextEnhanced struct {
	context.Context
	WorkspaceId     string
	TraceId         string
	parentVariables *KVStore
	variables       *KVStore
	Route           map[string]struct{}
	mu              sync.RWMutex
	IsTrace         bool
	startTime       time.Time
	Expiration      time.Time
	Extra           ContextExtra
	cancel          context.CancelFunc

	// 增强的资源管理
	resourceManager *ResourceManager
	referenceCount  int64
	cleaned         int64 // 使用atomic操作的清理标志
}

// ResourceManager 资源管理器
type ResourceManager struct {
	cleanupFuncs []func()
	mu           sync.Mutex
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

var EndParameters = "_end_"
var GenesisParameters = "_zero_"

type ExecutionStatus string

const (
	StatusPending   ExecutionStatus = "pending"
	StatusRunning   ExecutionStatus = "running"
	StatusPaused    ExecutionStatus = "paused"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
	StatusCanceled  ExecutionStatus = "canceled"
)

// NewResourceManager 创建资源管理器
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		cleanupFuncs: make([]func(), 0),
	}
}

// AddCleanupFunc 添加清理函数
func (rm *ResourceManager) AddCleanupFunc(cleanup func()) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.cleanupFuncs = append(rm.cleanupFuncs, cleanup)
}

// Cleanup 执行所有清理函数
func (rm *ResourceManager) Cleanup() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for i := len(rm.cleanupFuncs) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logx.Errorf("[Context] Cleanup function panic: %v", r)
				}
			}()
			rm.cleanupFuncs[i]()
		}()
	}

	rm.cleanupFuncs = rm.cleanupFuncs[:0] // 清空切片但保留底层数组
}

// Reset 重置执行上下文以便重用
func (ctx *ExecutionContextEnhanced) Reset() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	// 清理现有资源
	if ctx.resourceManager != nil {
		ctx.resourceManager.Cleanup()
	}

	// 重置变量存储
	if ctx.variables != nil {
		ctx.variables.Reset()
	}

	// 重置路由
	for k := range ctx.Route {
		delete(ctx.Route, k)
	}

	// 重置标志
	atomic.StoreInt64(&ctx.cleaned, 0)
	atomic.StoreInt64(&ctx.referenceCount, 1)

	// 重置其他字段
	ctx.WorkspaceId = ""
	ctx.TraceId = ""
	ctx.parentVariables = nil
	ctx.IsTrace = false
	ctx.startTime = time.Time{}
	ctx.Expiration = time.Time{}
	ctx.Extra = ContextExtra{}
	ctx.cancel = nil

	// 重新初始化资源管理器
	ctx.resourceManager = NewResourceManager()
}

// Initialize 初始化执行上下文
func (ctx *ExecutionContextEnhanced) Initialize(
	baseCtx context.Context,
	workspaceId string,
	traceId string,
	params map[string]any,
	extra ContextExtra,
	timeout time.Duration,
	nodeNum int64,
) {
	// 检查traceId是否以trace-开头
	if strings.HasPrefix(traceId, constants.TracePrefix) {
		ctx.IsTrace = true
	}

	ctx.Extra = extra
	withTimeout, cancel := context.WithTimeout(baseCtx, timeout)
	ctx.Context = withTimeout
	ctx.cancel = cancel
	ctx.Expiration = time.Now().Add(timeout)
	ctx.TraceId = traceId
	ctx.WorkspaceId = workspaceId

	// 初始化Route映射
	if ctx.Route == nil {
		ctx.Route = make(map[string]struct{})
	}

	ctx.startTime = time.Now()

	// 设置变量缓存
	if ctx.variables == nil {
		ctx.variables = CreateCache()
	}

	// 处理子上下文
	if extra.IsSub {
		if ect, ok := baseCtx.(*ExecutionContextEnhanced); ok {
			ctx.parentVariables = ect.variables
		} else {
			logx.Errorf("[Context] Failed to reuse variables: invalid context type")
		}
	}

	// 设置初始参数
	ctx.SetVariable(constants.GenesisParameters+workspaceId, params)
	ctx.Extra.NodeNum = nodeNum

	// 添加上下文取消的清理函数
	ctx.resourceManager.AddCleanupFunc(func() {
		if ctx.cancel != nil {
			ctx.cancel()
		}
	})
}

// AddRef 增加引用计数
func (ctx *ExecutionContextEnhanced) AddRef() {
	atomic.AddInt64(&ctx.referenceCount, 1)
}

// Release 减少引用计数，当引用计数为0时执行清理
func (ctx *ExecutionContextEnhanced) Release() {
	if atomic.AddInt64(&ctx.referenceCount, -1) == 0 {
		ctx.Cleanup()
	}
}

// Cleanup 清理资源
func (ctx *ExecutionContextEnhanced) Cleanup() {
	// 使用原子操作确保只清理一次
	if !atomic.CompareAndSwapInt64(&ctx.cleaned, 0, 1) {
		return
	}

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	// 执行资源清理
	if ctx.resourceManager != nil {
		ctx.resourceManager.Cleanup()
	}

	// 释放变量存储
	if ctx.variables != nil {
		ctx.variables.Recover()
		ctx.variables = nil
	}

	// 清理路由映射
	for k := range ctx.Route {
		delete(ctx.Route, k)
	}
}

// GetVariable 获取变量值（线程安全版本）
func (ctx *ExecutionContextEnhanced) GetVariable(key string) (map[string]any, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if ctx.variables != nil {
		if value, ok := ctx.variables.Get(key); ok {
			return value, true
		}
	}

	if ctx.parentVariables != nil {
		return ctx.parentVariables.Get(key)
	}

	return nil, false
}

// SetVariable 设置变量值（线程安全版本）
func (ctx *ExecutionContextEnhanced) SetVariable(key string, value map[string]any) {
	ctx.mu.RLock()
	variables := ctx.variables
	ctx.mu.RUnlock()

	if variables != nil {
		variables.Set(key, value)
	}
}

// SetRoute 设置路由（线程安全版本）
func (ctx *ExecutionContextEnhanced) SetRoute(nodeID string, route []string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	for _, r := range route {
		ctx.Route[nodeID+"_"+r] = struct{}{}
	}
}

// CheckRoute 检查路由（线程安全版本）
func (ctx *ExecutionContextEnhanced) CheckRoute(key []string) bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	for _, route := range key {
		if _, ok := ctx.Route[route]; ok {
			return true
		}
	}
	return false
}

// Cancel 取消执行
func (ctx *ExecutionContextEnhanced) Cancel() {
	if ctx.cancel != nil {
		ctx.cancel()
	}
}

// IsExpired 检查是否已过期
func (ctx *ExecutionContextEnhanced) IsExpired() bool {
	return time.Now().After(ctx.Expiration)
}

// GetReferenceCount 获取引用计数
func (ctx *ExecutionContextEnhanced) GetReferenceCount() int64 {
	return atomic.LoadInt64(&ctx.referenceCount)
}

// IsCleaned 检查是否已清理
func (ctx *ExecutionContextEnhanced) IsCleaned() bool {
	return atomic.LoadInt64(&ctx.cleaned) == 1
}

// GetExecutionDuration 获取执行时长
func (ctx *ExecutionContextEnhanced) GetExecutionDuration() time.Duration {
	return time.Since(ctx.startTime)
}

// AddCleanupFunc 添加清理函数
func (ctx *ExecutionContextEnhanced) AddCleanupFunc(cleanup func()) {
	if ctx.resourceManager != nil {
		ctx.resourceManager.AddCleanupFunc(cleanup)
	}
}

// Value 实现 context.Context 接口
func (ctx *ExecutionContextEnhanced) Value(key interface{}) interface{} {
	if key == "execution_context" {
		return ctx
	}
	return ctx.Context.Value(key)
}

// 增强的上下文池
var enhancedContextPool = sync.Pool{
	New: func() interface{} {
		ctx := &ExecutionContextEnhanced{
			Route:           make(map[string]struct{}),
			resourceManager: NewResourceManager(),
		}
		atomic.StoreInt64(&ctx.referenceCount, 1)
		return ctx
	},
}

// NewExecutionContextEnhanced 创建新的增强执行上下文
func NewExecutionContextEnhanced(
	ctx context.Context,
	workspaceId string,
	traceId string,
	params map[string]any,
	extra ContextExtra,
	timeout time.Duration,
	nodeNum int64,
) *ExecutionContextEnhanced {
	execCtx := enhancedContextPool.Get().(*ExecutionContextEnhanced)
	execCtx.Reset()
	execCtx.Initialize(ctx, workspaceId, traceId, params, extra, timeout, nodeNum)
	return execCtx
}

// ReleaseExecutionContextEnhanced 释放增强执行上下文
func ReleaseExecutionContextEnhanced(ctx *ExecutionContextEnhanced) {
	if ctx != nil {
		ctx.Release()
		enhancedContextPool.Put(ctx)
	}
}
