package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/components"
	"workflow/pkg/constants"
	"workflow/pkg/core"
	"workflow/pkg/errors"
)

// ExecutionManager 执行管理器，避免循环依赖
type ExecutionManager struct {
	poolManager *PoolManager
}

// PoolManager 简化的协程池管理器
type PoolManager struct {
	pool chan func()
	wg   sync.WaitGroup
}

// NewPoolManager 创建协程池管理器
func NewPoolManager(size int) *PoolManager {
	pm := &PoolManager{
		pool: make(chan func(), size),
	}

	// 启动worker
	for i := 0; i < size; i++ {
		go pm.worker()
	}

	return pm
}

// worker 工作协程
func (pm *PoolManager) worker() {
	for task := range pm.pool {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logx.Errorf("[Pool] Task panic: %v", r)
				}
				pm.wg.Done()
			}()
			task()
		}()
	}
}

// Submit 提交任务
func (pm *PoolManager) Submit(task func()) error {
	pm.wg.Add(1)
	select {
	case pm.pool <- task:
		return nil
	default:
		pm.wg.Done()
		return errors.ExecutionError("pool", "pool is full").Build()
	}
}

// Wait 等待所有任务完成
func (pm *PoolManager) Wait() {
	pm.wg.Wait()
}

// Close 关闭协程池
func (pm *PoolManager) Close() {
	close(pm.pool)
}

// NewExecutionManager 创建执行管理器
func NewExecutionManager(poolSize int) *ExecutionManager {
	return &ExecutionManager{
		poolManager: NewPoolManager(poolSize),
	}
}

// ExecutePhaseSimplified 简化的阶段执行
func (em *ExecutionManager) ExecutePhaseSimplified(execCtx *core.ExecutionContextEnhanced, runnable *Runnable, phase ExecutionPhase) error {
	logx.Debugf("[Workflow] Execute phase with %d nodes", len(phase.Nodes))
	// 创建错误收集器
	var errorCollector struct {
		mu     sync.Mutex
		errors []error
	}

	// 并发执行节点
	for _, node := range phase.Nodes {
		err := em.poolManager.Submit(func() {
			err := em.executeNodeSimplified(execCtx, runnable, node)
			if err != nil {
				errorCollector.mu.Lock()
				errorCollector.errors = append(errorCollector.errors, err)
				errorCollector.mu.Unlock()
			}
		})

		if err != nil {
			errorCollector.mu.Lock()
			errorCollector.errors = append(errorCollector.errors, err)
			errorCollector.mu.Unlock()
		}
	}

	// 等待所有任务完成
	em.poolManager.Wait()

	// 检查错误
	errorCollector.mu.Lock()
	defer errorCollector.mu.Unlock()

	if len(errorCollector.errors) > 0 {
		return fmt.Errorf("phase execution failed with %d errors: %v", len(errorCollector.errors), errorCollector.errors[0])
	}

	return nil
}

// executeNodeSimplified 简化的节点执行
func (em *ExecutionManager) executeNodeSimplified(execCtx *core.ExecutionContextEnhanced, runnable *Runnable, node *WorkflowNode) error {
	// 跳过路由检查对于Start组件
	if node.Type != constants.ComponentStart {
		if !em.shouldExecuteNode(execCtx, runnable, node) {
			logx.Debugf("[Workflow] Node route check failed, skipping: %s", node.ID)
			return em.handleSkippedNode(execCtx, node)
		}
	}

	// 创建组件
	component, err := components.ComponentFactory(nil, node.Type, node.Data)
	if err != nil {
		return fmt.Errorf("failed to create component [%s]: %w", node.ID, err)
	}
	defer component.Clear()

	// 执行节点
	return em.executeNodeSteps(execCtx, node, component)
}

// shouldExecuteNode 检查节点是否应该执行
func (em *ExecutionManager) shouldExecuteNode(execCtx *core.ExecutionContextEnhanced, runnable *Runnable, node *WorkflowNode) bool {
	route, ok := runnable.conditionRouter[node.ID]
	if !ok {
		return true
	}
	return execCtx.CheckRoute(route)
}

// handleSkippedNode 处理跳过的节点
func (em *ExecutionManager) handleSkippedNode(execCtx *core.ExecutionContextEnhanced, node *WorkflowNode) error {
	logx.Debugf("[Workflow] Node %s skipped, filling default values", node.ID)

	output := make(map[string]any)
	processedOutput, err := core.ProcessNodeOutput(output, node.Data.NodeDataOutputs)
	if err != nil {
		return err
	}

	execCtx.SetVariable(node.ID+".output", processedOutput)
	return nil
}

// executeNodeSteps 执行节点步骤
func (em *ExecutionManager) executeNodeSteps(execCtx *core.ExecutionContextEnhanced, node *WorkflowNode, component components.Component) error {
	logx.Debugf("[Workflow] Executing node: %s", node.ID)

	startTime := time.Now()
	var traceRecord *TraceRecore

	// 创建追踪记录（如果启用追踪）
	if execCtx.IsTrace && trace != nil {
		traceRecord = &TraceRecore{
			WorkspaceId: execCtx.WorkspaceId,
			TraceId:     execCtx.TraceId,
			NodeId:      node.ID,
			NodeName:    node.Data.Title,
			NodeType:    node.Type,
			Status:      "running",
			StartTime:   startTime,
			Step:        0,
		}
	}

	// 验证组件
	if validationErrors := component.Validate(); len(validationErrors) > 0 {
		return fmt.Errorf("component [%s] validation failed: %+v", node.ID, validationErrors)
	}

	// 准备输入
	input, err := em.prepareNodeInput(execCtx, node, component)
	if err != nil {
		// 记录错误并更新追踪
		if traceRecord != nil {
			traceRecord.Status = "failed"
			traceRecord.ErrorMsg = err.Error()
			traceRecord.ElapsedTime = time.Since(startTime).Milliseconds()
			trace.CreateTrace(execCtx, traceRecord)
		}
		return err
	}

	// 记录输入到追踪
	if traceRecord != nil {
		traceRecord.Input = input
		traceId, _ := trace.CreateTrace(execCtx, traceRecord)
		traceRecord.Id = traceId
	}

	// 执行组件
	result, err := component.Execute(execCtx, input)
	if err != nil {
		// 记录执行错误并更新追踪
		if traceRecord != nil {
			traceRecord.Status = "failed"
			traceRecord.ErrorMsg = err.Error()
			traceRecord.ElapsedTime = time.Since(startTime).Milliseconds()
			trace.UpdateTraceById(execCtx, traceRecord)
		}
		return err
	}

	// 处理输出
	output, err := em.processNodeOutput(input, result, node)
	if err != nil {
		// 记录输出处理错误并更新追踪
		if traceRecord != nil {
			traceRecord.Status = "failed"
			traceRecord.ErrorMsg = err.Error()
			traceRecord.ElapsedTime = time.Since(startTime).Milliseconds()
			trace.UpdateTraceById(execCtx, traceRecord)
		}
		return err
	}

	// 设置变量
	if node.Type == constants.ComponentEnd {
		execCtx.SetVariable(constants.EndParameters, output)
	} else {
		execCtx.SetVariable(node.ID+".output", output)
	}

	// 更新路由
	if result != nil {
		execCtx.SetRoute(node.ID, result.Route)
	}

	// 记录成功完成的追踪
	if traceRecord != nil {
		traceRecord.Status = "completed"
		traceRecord.Output = output
		traceRecord.ElapsedTime = time.Since(startTime).Milliseconds()
		trace.UpdateTraceById(execCtx, traceRecord)
	}

	logx.Debugf("[Workflow] Node %s executed successfully", node.ID)
	return nil
}

// prepareNodeInput 准备节点输入
func (em *ExecutionManager) prepareNodeInput(execCtx *core.ExecutionContextEnhanced, node *WorkflowNode, component components.Component) (map[string]any, error) {
	if node.Type == constants.ComponentStart {
		zero, ok := execCtx.GetVariable(constants.GenesisParameters + execCtx.WorkspaceId)
		if !ok {
			return nil, fmt.Errorf("input parameter not found: %s", constants.GenesisParameters+execCtx.WorkspaceId)
		}
		return core.ProcessNodeOutput(zero, node.Data.NodeDataOutputs)
	}

	if node.Type == constants.ComponentEnd {
		node.Data.NodeDataInputs.Properties = node.Data.NodeDataOutputs.Properties
		node.Data.NodeDataInputs.Required = node.Data.NodeDataOutputs.Required
		node.Data.NodeDataInputs.Type = node.Data.NodeDataOutputs.Type
	}

	mergedInput := make(map[string]any)

	// 获取自定义输入
	customInput, err := component.AnalyzeInputs(execCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze custom input: %w", err)
	}

	if customInput != nil {
		if customMap, ok := customInput.(map[string]any); ok {
			for k, v := range customMap {
				mergedInput[k] = v
			}
		}
	}

	// 获取标准输入
	standardInput, err := core.ParseNodeInputs(execCtx, node.Data.NodeDataInputsValue, node.Data.NodeDataInputs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse standard input: %w", err)
	}

	for k, v := range standardInput {
		mergedInput[k] = v
	}

	if len(mergedInput) == 0 {
		return map[string]any{}, nil
	}

	return mergedInput, nil
}

// processNodeOutput 处理节点输出
func (em *ExecutionManager) processNodeOutput(input map[string]any, result *core.Result, node *WorkflowNode) (map[string]any, error) {
	if node.Type == constants.ComponentEnd || node.Type == constants.ComponentEndItem {
		return input, nil
	}

	if result.Output == nil {
		return map[string]any{}, nil
	}

	return core.ProcessNodeOutput(result.Output.(map[string]any), node.Data.NodeDataOutputs)
}

// Close 关闭执行管理器
func (em *ExecutionManager) Close() {
	if em.poolManager != nil {
		em.poolManager.Close()
	}
}
