package engine

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rotisserie/eris"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/components"
	"workflow/pkg/core"
)

// ExecuteWorkflow 执行工作流
func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID string, serialID string, params map[string]any) error {
	// 参数校验
	if params == nil {
		return eris.New("params is nil")
	}

	logx.Debugf("[Workflow] Execute parameters: %+v", params)

	// 获取工作流执行器
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()

	if !ok {
		return eris.New("workflow not found: " + workflowID)
	}

	// 检查工作流状态
	if executor.status == core.WorkflowStatusShutdown {
		return eris.New("workflow is shutting down: " + workflowID)
	}

	// 更新访问时间
	atomic.StoreInt64(&executor.lastAccessed, time.Now().UnixNano())

	// 创建执行上下文
	executionContext := core.NewExecutionContext(ctx, workflowID, serialID, executor.totalNodes, params)
	// 创建带超时的上下文
	execCtx, cancel := context.WithTimeout(ctx, e.config.ExecutionTimeout)
	executionContext.Context = execCtx
	executionContext.Expiration = time.Now().Add(executor.defaultTTL)

	// 存储取消函数到执行上下文中，以便在暂停时可以取消
	executionContext.SetVariable("cancel", cancel)

	// 设置初始状态为运行中
	executionContext.State.Status = core.StatusRunning

	// 存储执行上下文（使用sync.Map避免锁争用）
	executor.execContexts.Store(serialID, executionContext)

	// 执行工作流
	err := e.executeWorkflowPhases(execCtx, executor, executionContext)

	// 处理执行结果
	if err != nil {
		executionContext.State.Status = core.StatusFailed
		executionContext.SetError("workflow", err)
		// 确保调用取消函数
		cancel()
		return err
	}

	executionContext.State.Status = core.StatusCompleted
	// 确保调用取消函数
	cancel()
	return nil
}

func (e *WorkflowEngine) ExecuteSingleWorkflow(ctx context.Context, workflowID, serialID string, nodeId string, params map[string]any) (*core.NodeResult, error) {
	startTime := time.Now()
	executor, ok := e.executorPool[workflowID]
	if !ok {
		return nil, eris.New("workflow not found: " + workflowID)
	}

	var node *core.Nodes
	// 遍历获取节点
	for _, n := range executor.definition.Nodes {
		if n.ID == nodeId {
			node = &n
			break
		}
	}

	nodeResult := &core.NodeResult{
		Input:    params,
		Route:    []string{components.Success},
		NodeID:   nodeId,
		Duration: time.Since(startTime).Milliseconds(),
		Type:     node.Type,
		NodeName: node.Data.Title,
	}
	// 创建执行上下文
	execCtx := core.NewExecutionContext(ctx, workflowID, serialID, executor.totalNodes, params)
	execCtx.Context = ctx
	execCtx.Expiration = time.Now().Add(executor.defaultTTL)

	// 执行节点
	component, err := components.ComponentFactory(e, node.Type, node.Data)
	if err != nil {
		logx.Errorw("failed to create component", logx.Field("error", err.Error()))
		nodeResult.Error = err.Error()
		return nodeResult, eris.New("failed to create component: " + err.Error())
	}
	result, err := e.executeNode(execCtx, 0, nodeId, node, component, params)
	if err != nil {
		logx.Errorw("failed to execute node", logx.Field("error", err.Error()))
		nodeResult.Error = err.Error()
		return nodeResult, eris.New("failed to execute node: " + err.Error())
	}
	nodeResult.Output = result.Output
	nodeResult.Duration = time.Since(startTime).Milliseconds()
	return nodeResult, nil
}

// GetNodeResult 获取指定工作流节点的执行结果
func (e *WorkflowEngine) GetNodeResult(workflowID, serialID, nodeID string) (*core.NodeResult, bool) {
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()

	if !ok {
		logx.Debugf("[Workflow] Get node result, workflow not found: %s", workflowID)
		return nil, false
	}

	// 使用sync.Map获取执行上下文
	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		logx.Debugf("[Workflow] Get node result, execution context not found: %s", serialID)
		return nil, false
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		logx.Debugf("[Workflow] Get node result, execution context type error: %s", serialID)
		return nil, false
	}

	return execCtx.GetNodeResult(nodeID)
}

// executeWorkflowPhases 执行工作流阶段
func (e *WorkflowEngine) executeWorkflowPhases(ctx context.Context, executor *Executor, execCtx *core.ExecutionContext) error {
	// 总执行计划
	logx.Debugf("[Workflow] Total execution plan: %d", len(executor.executionPlan.Phases))
	for phaseIdx, phase := range executor.executionPlan.Phases {
		if err := e.executePhase(ctx, executor, execCtx, phase, phaseIdx); err != nil {
			logx.Errorw("[Workflow] Execute phase failed", logx.Field("phase index", phaseIdx), logx.Field("error", err.Error()))
		}
	}
	return nil
}

// executePhase 执行单个阶段
func (e *WorkflowEngine) executePhase(ctx context.Context, executor *Executor, execCtx *core.ExecutionContext, phase ExecutionPhase, phaseIdx int) error {
	// 添加上下文超时控制
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	logx.Debugf("[Workflow] Execute phase: %d, node count: %d", phaseIdx, len(phase.Nodes))
	// 获取单例协程池
	var execWg sync.WaitGroup
	pool := GetGlobalPool()

	// 创建错误收集器
	var errMu sync.Mutex
	var errorStack []error

	for i, node := range phase.Nodes {
		execWg.Add(1)
		err := pool.Submit(func() {
			defer execWg.Done()
			err := e.handleNodeExecution(execCtx, phaseIdx*1000+i, executor, node)
			if err != nil {
				errField := eris.ToString(err, true)
				logx.Errorf("[Workflow] Execute node failed, id:%s, name:%s, type:%s, traceId:%s, error %s",
					node.ID, node.Data.Title, node.Type, execCtx.TraceId, errField)
				// 收集错误
				errMu.Lock()
				errorStack = append(errorStack, fmt.Errorf("node[%s] execute failed: %w", node.ID, err))
				errMu.Unlock()
			}
		})
		if err != nil {
			logx.Errorf("[Workflow] Submit node failed: id:%s, name:%s, type:%s, traceId:%s, error:%s",
				node.ID, node.Data.Title, node.Type, execCtx.TraceId, err.Error())
			// 收集提交错误
			errMu.Lock()
			errorStack = append(errorStack, fmt.Errorf("节点[%s]提交失败: %w", node.ID, err))
			errMu.Unlock()
		}
	}
	execWg.Wait()

	// 如果有错误，返回组合错误
	if len(errorStack) > 0 {
		return fmt.Errorf("phase[%d] execute failed: %v", phaseIdx, errorStack)
	}
	return nil
}

// handleNodeExecution 处理节点执行
func (e *WorkflowEngine) handleNodeExecution(execCtx *core.ExecutionContext, phaseIdx int, executor *Executor, node *WorkflowNode) error {
	// 开始类型的组件跳过路由检查
	if node.Type == components.Start || node.Type == components.StartItem {
		return e.executeStartNode(execCtx, int64(phaseIdx), node)
	}

	// 检查路由
	ok, err := e.checkNodeRoute(execCtx, executor, node)
	if err != nil {
		return err
	}
	if !ok {
		logx.Debugf("[Workflow] Node route skipped: %s", node.ID)
		return nil
	}

	// 创建并执行组件
	component, err := components.ComponentFactory(e, node.Type, node.Data)
	if err != nil {
		return eris.New("failed to create component [" + node.ID + "]: " + err.Error())
	}

	// 直接执行节点，不使用线程池
	_, err = e.executeNode(execCtx, int64(phaseIdx), node.ID, node.Nodes, component, nil)
	if err != nil {
		return err
	}

	return nil
}

// executeStartNode 执行开始节点
func (e *WorkflowEngine) executeStartNode(execCtx *core.ExecutionContext, phaseIdx int64, node *WorkflowNode) error {
	params, bool := execCtx.GetVariable("_zero")
	if !bool {
		return eris.New("no data input")
	}
	output, err := core.ProcessNodeOutput(params.(map[string]any), node.Data.NodeDataOutputs)
	var errorMsg string
	if err != nil {
		output = map[string]any{}
		errorMsg = err.Error()
	}
	if execCtx.IsTrace {
		// 创建 trace
		Trace.CreateTrace(execCtx, &TraceRecore{
			WorkspaceId: execCtx.WorkspaceId,
			TraceId:     execCtx.TraceId,
			NodeId:      node.ID,
			NodeName:    node.Data.Title,
			NodeType:    node.Type,
			Logic:       "",
			Input:       params,
			Output:      output,
			Step:        phaseIdx,
			Status:      string(core.StatusCompleted),
			StartTime:   time.Now(),
			ElapsedTime: 0,
			ErrorMsg:    errorMsg,
		})
	}

	execCtx.SetVariable(node.ID+".output", output)
	nodeResult := &core.NodeResult{
		Input:    params,
		Output:   output,
		Route:    []string{components.Success},
		NodeName: node.Data.Title,
		NodeID:   node.ID,
		Duration: 0,
		Error:    errorMsg,
		Type:     node.Type,
	}
	_ = e.updateWorkflowState(execCtx, nil, []*core.NodeResult{nodeResult})
	return err
}

// checkNodeRoute 检查节点路由
func (e *WorkflowEngine) checkNodeRoute(execCtx *core.ExecutionContext, executor *Executor, node *WorkflowNode) (bool, error) {
	route, ok := executor.conditionRouter[node.ID]
	if !ok {
		logx.Debugf("[Workflow] Node route not found, skip[loop single node]: %s", node.ID)
		return true, nil
	}

	if !execCtx.CheckRoute(route) {
		logx.Debugf("路由:%+v\n", route)
		return false, e.handleSkippedNode(execCtx, node)
	}

	return true, nil
}

// handleSkippedNode 处理跳过的节点
func (e *WorkflowEngine) handleSkippedNode(execCtx *core.ExecutionContext, node *WorkflowNode) error {
	logx.Debugf("[Workflow] %s node necessary route not found, fill empty parameters\n", node.ID)
	// 补齐默认零值
	input := make(map[string]any)
	output, err := core.ProcessNodeOutput(input, node.Data.NodeDataOutputs)
	if err != nil {
		return err
	}

	execCtx.SetVariable(node.ID+".output", output)

	nodeResult := &core.NodeResult{
		Input:    input,
		Output:   output,
		Route:    []string{components.Skip},
		NodeID:   node.ID,
		Duration: 0,
		Error:    "",
		Type:     node.Type,
	}

	return e.updateWorkflowState(execCtx, nil, []*core.NodeResult{nodeResult})
}

// prepareNodeInput 准备节点输入数据
func (e *WorkflowEngine) prepareNodeInput(ctx *core.ExecutionContext, err error, node *core.Nodes, component components.Component) (any, error) {
	if err != nil {
		return nil, err
	}

	if node.Type == "start" {
		zero, ok := ctx.GetVariable("_zero")
		if !ok {
			return nil, eris.New("input parameter not found")
		}
		return core.ProcessNodeOutput(zero.(map[string]any), node.Data.NodeDataOutputs)
	}
	if node.Type == "end" {
		// 结束节点 输入和输出相同,DSL 定义了输出, 所以需要将输出赋值给输入
		node.Data.NodeDataInputs.Properties = node.Data.NodeDataOutputs.Properties
		node.Data.NodeDataInputs.Required = node.Data.NodeDataOutputs.Required
		node.Data.NodeDataInputs.Type = node.Data.NodeDataOutputs.Type
	}

	// 初始化合并结果
	mergedInput := make(map[string]any)

	// 1. 获取自定义输入
	customInput, err := component.AnalyzeInputs(ctx)
	if err != nil {
		return nil, eris.New("failed to analyze custom input: " + err.Error())
	}

	// 2. 如果有自定义输入，先复制到合并结果中
	if customInput != nil {
		if customMap, ok := customInput.(map[string]any); ok {
			for k, v := range customMap {
				mergedInput[k] = v
			}
		} else {
			return nil, eris.New("custom input must be map type")
		}
	}

	// 3. 获取标准输入
	standardInput, err := core.ParseNodeInputs(ctx, node.Data.NodeDataInputsValue, node.Data.NodeDataInputs)
	if err != nil {
		return nil, eris.New("failed to parse standard input: " + err.Error())
	}

	// 4. 合并标准输入
	maps.Copy(mergedInput, standardInput)

	// 5. 如果没有任何输入，返回空map
	if len(mergedInput) == 0 {
		return map[string]any{}, nil
	}

	return mergedInput, nil
}

// processNodeOutput 处理节点输出数据
func (e *WorkflowEngine) processNodeOutput(input any, err error, result *core.Result, node *core.Nodes) (any, error) {
	if err != nil {
		return nil, err
	}

	if node.Type == components.End || node.Type == components.EndItem {
		return input, nil
	}

	if result.Output == nil {
		return map[string]any{}, nil
	}

	output, err := core.ProcessNodeOutput(result.Output.(map[string]any), node.Data.NodeDataOutputs)
	if err != nil {
		return nil, eris.New("failed to process output mapping: " + err.Error())
	}
	return output, nil
}

// executeNode 执行节点
func (e *WorkflowEngine) executeNode(ctx *core.ExecutionContext, step int64, nodeID string, node *core.Nodes, component components.Component, singleParam map[string]any) (*core.Result, error) {
	startTime := time.Now()

	// 1. 验证组件
	if err := e.validateComponent(component, nodeID); err != nil {
		return nil, err
	}

	// 2. 准备输入数据
	var input any
	var err error
	if singleParam != nil {
		input = singleParam
	} else {
		if input, err = e.prepareNodeInput(ctx, nil, node, component); err != nil {
			return nil, err
		}
	}

	// 创建 trace
	if ctx.IsTrace {
		trace := &TraceRecore{
			WorkspaceId: ctx.WorkspaceId,
			TraceId:     ctx.TraceId,
			NodeId:      nodeID,
			NodeName:    node.Data.Title,
			NodeType:    node.Type,
			Input:       input,
			Logic:       "",
			StartTime:   startTime,
			Step:        step,
		}
		Trace.CreateTrace(ctx, trace)
	}

	// 3. 执行组件
	result, err := e.executeComponent(ctx, nil, component, input)
	if err != nil {
		// 更新 trace 记录错误信息
		if ctx.IsTrace {
			Trace.UpdateTrace(ctx, &TraceRecore{
				NodeId:      nodeID,
				TraceId:     ctx.TraceId,
				Output:      nil,
				Status:      string(core.StatusFailed),
				ElapsedTime: time.Since(startTime).Milliseconds(),
				ErrorMsg:    err.Error(),
			})
		}
		return nil, err
	}

	// 4. 处理输出数据
	output, err := e.processNodeOutput(input, nil, result, node)
	if err != nil {
		// 更新 trace 记录错误信息
		if ctx.IsTrace {
			Trace.UpdateTrace(ctx, &TraceRecore{
				NodeId:      nodeID,
				TraceId:     ctx.TraceId,
				Output:      nil,
				Status:      string(core.StatusFailed),
				ElapsedTime: time.Since(startTime).Milliseconds(),
				ErrorMsg:    err.Error(),
			})
		}
		return nil, err
	}

	// 5. 释放组件
	e.Clear(component)

	// 更新 trace
	if ctx.IsTrace {
		Trace.UpdateTrace(ctx, &TraceRecore{
			NodeId:      nodeID,
			TraceId:     ctx.TraceId,
			Output:      output,
			Status:      string(core.StatusCompleted),
			ElapsedTime: time.Since(startTime).Milliseconds(),
			ErrorMsg:    "",
		})
	}

	// 5. 更新上下文
	if err := e.updateNodeContext(ctx, nil, nodeID, output); err != nil {
		return nil, e.handleNodeError(ctx, node, err, startTime)
	}

	// 6. 创建并返回结果
	nodeResult := e.createNodeResult(node, input, output, result, startTime)
	if err := e.updateWorkflowState(ctx, nil, []*core.NodeResult{nodeResult}); err != nil {
		logx.Errorf("[Workflow] Update failed [WorkflowID:%s] [SerialID:%s] [NodeID:%s] [Error:%v]",
			ctx.WorkspaceId, ctx.TraceId, node.ID, err)
		return nil, err
	}
	outputJson, err := sonic.Marshal(output)
	if err != nil {
		logx.Errorf("[Workflow] Output result serialization failed [WorkflowID:%s] [SerialID:%s] [NodeID:%s] [Error:%v]",
			ctx.WorkspaceId, ctx.TraceId, node.ID, err)
		return nil, err
	}

	logx.Infof("[Workflow] Node execution completed traceId: %s, node id: %s, output: %s", ctx.TraceId, nodeID, string(outputJson))
	return result, nil
}

// validateComponent 验证组件
func (e *WorkflowEngine) validateComponent(component components.Component, nodeID string) error {
	validateErrors := component.Validate()
	exception := component.Exception()
	if !reflect.ValueOf(exception).IsZero() {
		validateErrors = exception.Validate()
	}
	if len(validateErrors) > 0 {
		return eris.New("component [" + nodeID + "] validation failed: " + fmt.Sprintf("%+v", validateErrors))
	}
	return nil
}

// executeComponent 执行组件
func (e *WorkflowEngine) executeComponent(ctx *core.ExecutionContext, err error, component components.Component, input any) (*core.Result, error) {
	if err != nil {
		return nil, err
	}

	exceptionCnf := component.Exception()
	if reflect.ValueOf(exceptionCnf).IsZero() {
		return component.Execute(ctx, input)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(exceptionCnf.Timeout)*time.Second)
	defer cancel()

	var result *core.Result
	for i := 0; i < exceptionCnf.RetryTimes; i++ {
		result, err = component.Execute(timeoutCtx, input)
		if err == nil && ctx.Err() == nil {
			break
		}
		logx.Debugf("[Workflow] Execute failed, retry: %s", err.Error())
	}
	return result, err
}

// updateNodeContext 更新节点上下文
func (e *WorkflowEngine) updateNodeContext(ctx *core.ExecutionContext, err error, nodeID string, output any) error {
	if err != nil {
		return err
	}
	if output != nil {
		logx.Debugf("[Workflow] Update node context: %s, %+v", nodeID, output)
		ctx.SetVariable(nodeID+".output", output)
	}
	return nil
}

// createNodeResult 创建节点结果
func (e *WorkflowEngine) createNodeResult(node *core.Nodes, input, output any, result *core.Result, startTime time.Time) *core.NodeResult {
	return &core.NodeResult{
		Input:    input,
		Output:   output,
		Route:    result.Route,
		NodeID:   node.ID,
		Duration: time.Since(startTime).Milliseconds(),
		Error:    "",
		Type:     node.Type,
	}
}

// handleNodeError 处理节点错误
func (e *WorkflowEngine) handleNodeError(ctx *core.ExecutionContext, node *core.Nodes, err error, startTime time.Time) error {
	nodeResult := &core.NodeResult{
		NodeID:   node.ID,
		Duration: time.Since(startTime).Milliseconds(),
		Error:    err.Error(),
		Type:     node.Type,
	}
	e.updateWorkflowState(ctx, err, []*core.NodeResult{nodeResult})
	return err
}

// updateWorkflowState 更新工作流状态
func (e *WorkflowEngine) updateWorkflowState(ctx *core.ExecutionContext, err error, results []*core.NodeResult) error {
	if err != nil {
		return err
	}

	for _, result := range results {
		ctx.SetNodeResult(result.NodeID, result)
	}

	// 计算整体进度
	completed := len(ctx.State.Result)
	total := int(ctx.TotalNodes)
	ctx.State.Progress = float64(completed) / float64(total)

	// 使用安全的方法进行序列化
	resjson, err := ctx.MarshalResult()
	if err != nil {
		logx.Errorf("[Workflow] Serialization failed [WorkflowID:%s] [SerialID:%s] [Error:%v]",
			ctx.WorkspaceId, ctx.TraceId, err)
		return err
	}

	logx.Infof("[Workflow] Phase execution result traceId: %s, workflowId: %s, output: %s", ctx.TraceId, ctx.WorkspaceId, string(resjson))
	return nil
}
