package engine

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"

	"workflow/pkg/components"
	"workflow/pkg/core"
)

// WorkflowEngine 工作流引擎 https://deepwiki.com/XXueTu/workflow/1-overview
type WorkflowEngine struct {
	// 核心组件
	pool     *ants.PoolWithFunc // 任务执行池
	poolSize int                // 当前池大小
	// 配置
	config *EngineConfig // 引擎配置

	// 工作流管理
	executorPool map[string]*Executor // 工作流执行器池
	mu           sync.RWMutex         // 读写锁

	// 资源管理
	cleanup    chan string   // 清理通道
	shutdownCh chan struct{} // 关闭信号通道

}

// NewWorkflowEngine 创建新的工作流引擎
func NewWorkflowEngine(opts ...Option) *WorkflowEngine {
	config := DefaultConfig()
	engine := &WorkflowEngine{
		executorPool: make(map[string]*Executor),
		config:       config,
		cleanup:      make(chan string, 100),
		shutdownCh:   make(chan struct{}),
	}

	for _, opt := range opts {
		opt(engine)
	}

	// 创建任务处理函数
	taskHandler := func(i any) {
		task, ok := i.(*ExecutionTask)
		if !ok {
			return
		}
		defer task.Sw.Done()
		result, err := task.Engine.executeNode(task.Context, task.Step, task.NodeID, task.Executor.definition.Nodes[task.NodeID], task.Component, nil)
		if err != nil {
			logx.Errorf("[工作流执行] 节点执行失败 [工作流ID:%s] [序列ID:%s] [节点ID:%s] [步骤:%d] [错误:%v]",
				task.WorkflowID, task.SerialID, task.NodeID, task.Step, err)
			task.Context.SetError(task.NodeID, err)
			return
		}
		logx.Infof("[工作流执行] 节点执行成功 [工作流ID:%s] [序列ID:%s] [节点ID:%s] [步骤:%d] [输出:%+v]",
			task.WorkflowID, task.SerialID, task.NodeID, task.Step, result.Output)
	}

	// 创建 ants 函数池
	pool, err := ants.NewPoolWithFunc(engine.config.InitialPoolSize, taskHandler,
		ants.WithPreAlloc(true),
		ants.WithMaxBlockingTasks(engine.config.MaxConcurrentWorkflows),
		ants.WithNonblocking(true),
	)
	if err != nil {
		panic(err)
	}
	engine.pool = pool
	engine.poolSize = engine.config.InitialPoolSize

	// 启动清理和指标收集协程
	engine.startCleanupRoutine()

	return engine
}

// RegisterWorkflow 注册工作流
func (e *WorkflowEngine) RegisterWorkflow(ctx context.Context, def *core.WorkflowDef) error {

	// 构建执行计划
	plan, err := e.buildExecutionPlan(ctx, def)
	if err != nil {
		return errors.New("无法构建执行计划: " + err.Error())
	}

	// 构建条件路由
	condition, err := e.buildConditionRouter(def)
	if err != nil {
		return errors.New("无法构建条件映射: " + err.Error())
	}

	// 创建新的执行器
	executor := e.createExecutor(def, plan, condition)

	// 注册执行器
	e.registerExecutor(def.ID, executor)

	return nil
}

// createExecutor 创建执行器
func (e *WorkflowEngine) createExecutor(def *core.WorkflowDef, plan *ExecutionPlan, condition map[string][]string) *Executor {
	e.mu.Lock()
	defer e.mu.Unlock()
	executor := &Executor{
		definition:      def,
		env:             make(map[string]any),
		executionPlan:   plan,
		conditionRouter: condition,
		totalNodes:      int64(len(def.Nodes)),
		status:          core.WorkflowStatusActive,
		createdAt:       time.Now(),
		defaultTTL:      e.config.DefaultContextTTL,
		shutdownCh:      make(chan struct{}),
	}

	if existing, exists := e.executorPool[def.ID]; exists {
		atomic.StoreInt64(&existing.lastAccessed, time.Now().UnixNano())
		existing.status = core.WorkflowStatusDeploying
		// 复制已经存在的执行上下文,防止删除正在执行的上下文
		existing.execContexts.Range(func(key, value any) bool {
			executor.execContexts.Store(key, value)
			return true
		})
	}
	return executor
}

// registerExecutor 注册执行器
func (e *WorkflowEngine) registerExecutor(workflowID string, executor *Executor) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.executorPool[workflowID] = executor
}

// DeregisterWorkflow 注销工作流
func (e *WorkflowEngine) DeregisterWorkflow(ctx context.Context, workflowID string, graceful bool) error {
	e.mu.Lock()
	executor, exists := e.executorPool[workflowID]
	if !exists {
		e.mu.Unlock()
		return errors.New("工作流未找到: " + workflowID)
	}

	executor.status = core.WorkflowStatusShutdown
	e.mu.Unlock()

	if !graceful {
		return e.immediateDeregister(workflowID, executor)
	}

	return e.gracefulDeregister(ctx, workflowID, executor)
}

// immediateDeregister 立即注销工作流
func (e *WorkflowEngine) immediateDeregister(workflowID string, executor *Executor) error {
	e.mu.Lock()
	delete(e.executorPool, workflowID)
	e.mu.Unlock()
	close(executor.shutdownCh)
	return nil
}

// gracefulDeregister 优雅注销工作流
func (e *WorkflowEngine) gracefulDeregister(ctx context.Context, workflowID string, executor *Executor) error {
	go func() {
		timeout := time.After(30 * time.Second)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if atomic.LoadInt64(&executor.activeCount) == 0 {
					e.immediateDeregister(workflowID, executor)
					return
				}
			case <-timeout:
				e.immediateDeregister(workflowID, executor)
				return
			case <-ctx.Done():
				e.immediateDeregister(workflowID, executor)
				return
			}
		}
	}()

	return nil
}

// ExecuteWorkflow 执行工作流
func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID string, serialID string, params map[string]any) error {
	// 参数校验
	if params == nil {
		return errors.New("参数为空")
	}

	logx.Debugf("[工作流] 执行参数: %+v", params)

	// 获取工作流执行器
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()

	if !ok {
		return errors.New("工作流未找到: " + workflowID)
	}

	// 检查工作流状态
	if executor.status == core.WorkflowStatusShutdown {
		return errors.New("工作流正在关闭: " + workflowID)
	}

	// 更新访问时间
	atomic.StoreInt64(&executor.lastAccessed, time.Now().UnixNano())

	// 创建带超时的上下文
	execCtx, cancel := context.WithTimeout(ctx, e.config.ExecutionTimeout)

	// 创建执行上下文
	executionContext := core.NewExecutionContext(execCtx, workflowID, serialID, executor.totalNodes, params)
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
		return nil, errors.New("工作流未找到: " + workflowID)
	}
	node := executor.definition.Nodes[nodeId]
	nodeResult := &core.NodeResult{
		Input:    params,
		Route:    []string{components.Success},
		NodeID:   nodeId,
		Duration: time.Since(startTime).Milliseconds(),
		Type:     node.Type,
		NodeName: node.Name,
	}
	// 创建执行上下文
	execCtx := core.NewExecutionContext(ctx, workflowID, serialID, executor.totalNodes, params)
	execCtx.Context = ctx
	execCtx.Expiration = time.Now().Add(executor.defaultTTL)

	// 执行节点
	component, err := components.ComponentFactory(e, node.Type, node)
	if err != nil {
		logx.Errorw("无法创建组件", logx.Field("error", err.Error()))
		nodeResult.Error = err.Error()
		return nodeResult, errors.New("无法创建组件: " + err.Error())
	}
	result, err := e.executeNode(execCtx, 0, nodeId, node, component, params)
	if err != nil {
		logx.Errorw("执行节点失败", logx.Field("error", err.Error()))
		nodeResult.Error = err.Error()
		return nodeResult, errors.New("执行节点失败: " + err.Error())
	}
	nodeResult.Output = result.Output
	nodeResult.Duration = time.Since(startTime).Milliseconds()
	return nodeResult, nil
}

// executeWorkflowPhases 执行工作流阶段
func (e *WorkflowEngine) executeWorkflowPhases(ctx context.Context, executor *Executor, execCtx *core.ExecutionContext) error {
	// 总执行计划
	logx.Debugf("[工作流] 总执行计划: %d", len(executor.executionPlan.Phases))
	for phaseIdx, phase := range executor.executionPlan.Phases {
		if err := e.executePhase(ctx, executor, execCtx, phase, phaseIdx); err != nil {
			return err
		}
	}
	return nil
}

// executePhase 执行单个阶段
func (e *WorkflowEngine) executePhase(ctx context.Context, executor *Executor, execCtx *core.ExecutionContext, phase ExecutionPhase, phaseIdx int) error {
	var sw sync.WaitGroup
	nodes := make([]string, len(phase.Nodes))
	g, ctx := errgroup.WithContext(ctx)
	logx.Debugf("[工作流] 执行阶段: %d, 节点数: %d", phaseIdx, len(phase.Nodes))
	for i, node := range phase.Nodes {
		if err := e.submitNodeTask(ctx, phaseIdx*10000+i, executor, execCtx, node, &sw, g, nodes, i); err != nil {
			return err
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}
	sw.Wait()

	return nil
}

// submitNodeTask 提交节点任务
func (e *WorkflowEngine) submitNodeTask(ctx context.Context, phaseIdx int, executor *Executor, execCtx *core.ExecutionContext, node *WorkflowNode, sw *sync.WaitGroup, g *errgroup.Group, nodes []string, index int) error {
	logx.Debugf("[工作流] 执行节点: %s, 步骤: %d", node.ID, phaseIdx)
	nodes[index] = node.ID

	// 检查上下文状态
	if err := e.checkContext(ctx, executor); err != nil {
		return err
	}

	// 获取节点定义副本
	nodeCopy := node

	// 增加等待计数
	sw.Add(1)

	// 提交任务
	g.Go(func() error {
		if err := e.handleNodeExecution(execCtx, phaseIdx, executor, nodeCopy, sw); err != nil {
			return err
		}
		return nil
	})

	return nil
}

// checkContext 检查上下文状态
func (e *WorkflowEngine) checkContext(ctx context.Context, executor *Executor) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-executor.shutdownCh:
		return errors.New("工作流正在关闭")
	default:
		return nil
	}
}

// handleNodeExecution 处理节点执行
func (e *WorkflowEngine) handleNodeExecution(execCtx *core.ExecutionContext, phaseIdx int, executor *Executor, node *WorkflowNode, sw *sync.WaitGroup) error {
	// 开始类型的组件跳过路由检查
	if node.Type == components.Start || node.Type == components.StartItem {
		sw.Done()
		return e.executeStartNode(execCtx, int64(phaseIdx), node)
	}

	// 检查路由
	if err := e.checkNodeRoute(execCtx, executor, node); err != nil {
		return err
	}

	// 创建并执行组件
	component, err := components.ComponentFactory(e, node.Type, node.NodeDefinition)
	if err != nil {
		return errors.New("无法创建组件 [" + node.ID + "]: " + err.Error())
	}

	task := &ExecutionTask{
		NodeID:     node.ID,
		Component:  component,
		Context:    execCtx,
		Engine:     e,
		Executor:   executor,
		WorkflowID: executor.definition.ID,
		SerialID:   execCtx.TraceId,
		Sw:         sw,
		Step:       int64(phaseIdx),
	}

	return e.pool.Invoke(task)
}

// executeStartNode 执行开始节点
func (e *WorkflowEngine) executeStartNode(execCtx *core.ExecutionContext, phaseIdx int64, node *WorkflowNode) error {
	params, bool := execCtx.GetVariable("_zero")
	if !bool {
		return errors.New("无数据输入")
	}
	output, err := core.ProcessNodeOutput(params.(map[string]any), node.Outputs)
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
			NodeName:    node.Name,
			Logic:       node.NodeDefinition.Config,
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
		NodeName: node.Name,
		NodeID:   node.ID,
		Duration: 0,
		Error:    errorMsg,
		Type:     node.Type,
	}
	e.updateWorkflowState(execCtx, nil, []*core.NodeResult{nodeResult})
	return err
}

// checkNodeRoute 检查节点路由
func (e *WorkflowEngine) checkNodeRoute(execCtx *core.ExecutionContext, executor *Executor, node *WorkflowNode) error {
	route, ok := executor.conditionRouter[node.ID]
	if !ok {
		return errors.New("节点初始化路由未找到: " + node.ID)
	}

	if !execCtx.CheckRoute(route) {
		logx.Debugf("路由:%+v\n", route)
		return e.handleSkippedNode(execCtx, node)
	}

	return nil
}

// handleSkippedNode 处理跳过的节点
func (e *WorkflowEngine) handleSkippedNode(execCtx *core.ExecutionContext, node *WorkflowNode) error {
	logx.Debugf("%s节点必要路由未找到,跳过\n", node.ID)

	output, err := core.ProcessNodeOutput(map[string]any{}, node.Outputs)
	if err != nil {
		return err
	}

	execCtx.SetVariable(node.ID+".output", output)

	nodeResult := &core.NodeResult{
		Input:    map[string]any{},
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
func (e *WorkflowEngine) prepareNodeInput(ctx *core.ExecutionContext, err error, node *core.NodeDefinition, component components.Component) (any, error) {
	if err != nil {
		return nil, err
	}

	if node.Type == "start" {
		zero, ok := ctx.GetVariable("_zero")
		if !ok {
			return nil, errors.New("输入参数未找到")
		}
		return core.ProcessNodeOutput(zero.(map[string]any), node.Outputs)
	}

	// 初始化合并结果
	mergedInput := make(map[string]any)

	// 1. 获取自定义输入
	customInput, err := component.AnalyzeInputs(ctx)
	if err != nil {
		return nil, errors.New("分析自定义输入失败: " + err.Error())
	}

	// 2. 如果有自定义输入，先复制到合并结果中
	if customInput != nil {
		if customMap, ok := customInput.(map[string]any); ok {
			for k, v := range customMap {
				mergedInput[k] = v
			}
		} else {
			return nil, errors.New("自定义输入必须是map类型")
		}
	}

	// 3. 获取标准输入
	standardInput, err := core.ParseNodeInputs(node.Inputs, ctx)
	if err != nil {
		return nil, errors.New("解析标准输入失败: " + err.Error())
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
func (e *WorkflowEngine) processNodeOutput(input any, err error, result *core.Result, node *core.NodeDefinition) (any, error) {
	if err != nil {
		return nil, err
	}

	if node.Type == components.End || node.Type == components.EndItem {
		return input, nil
	}

	if result.Output == nil {
		return map[string]any{}, nil
	}

	output, err := core.ProcessNodeOutput(result.Output.(map[string]any), node.Outputs)
	if err != nil {
		return nil, errors.New("处理输出映射失败: " + err.Error())
	}
	return output, nil
}

// executeNode 执行节点
func (e *WorkflowEngine) executeNode(ctx *core.ExecutionContext, step int64, nodeID string, node *core.NodeDefinition, component components.Component, singleParam map[string]any) (*core.Result, error) {

	startTime := time.Now()

	// 1. 验证组件
	err := e.validateComponent(component, node)

	// 2. 准备输入数据
	var input any
	if singleParam != nil {
		input = singleParam
	} else {
		input, err = e.prepareNodeInput(ctx, err, node, component)
	}
	// 创建 trace
	trace := &TraceRecore{
		WorkspaceId: ctx.WorkspaceId,
		TraceId:     ctx.TraceId,
		NodeId:      nodeID,
		NodeName:    node.Name,
		Input:       input,
		Logic:       node.Config,
		StartTime:   startTime,
		Step:        step,
	}
	Trace.CreateTrace(ctx, trace)
	// 3. 执行组件
	result, err := e.executeComponent(ctx, err, component, input)

	// 4. 处理输出数据
	output, err := e.processNodeOutput(input, err, result, node)

	var errorMsg string
	var status = string(core.StatusCompleted)
	if err != nil {
		status = string(core.StatusFailed)
		errorMsg = err.Error()
		output = map[string]any{}
	}
	// 5. 更新 trace
	Trace.UpdateTrace(ctx, &TraceRecore{
		NodeId:      nodeID,
		TraceId:     ctx.TraceId,
		Output:      output,
		Status:      status,
		ElapsedTime: time.Since(startTime).Milliseconds(),
		ErrorMsg:    errorMsg,
	})

	// 5. 更新上下文
	err = e.updateNodeContext(ctx, err, nodeID, output)
	if err != nil {
		return nil, e.handleNodeError(ctx, node, err, startTime)
	}

	// 6. 创建并返回结果
	nodeResult := e.createNodeResult(node, input, output, result, startTime)
	err = e.updateWorkflowState(ctx, err, []*core.NodeResult{nodeResult})
	if err != nil {
		logx.Errorf("[工作流状态] 更新失败 [工作流ID:%s] [序列ID:%s] [节点ID:%s] [错误:%v]",
			ctx.WorkspaceId, ctx.TraceId, node.ID, err)
	}

	if err != nil {
		logx.Errorw("[引擎] 执行节点失败",
			logx.Field("节点ID", nodeID),
			logx.Field("错误", err))
		return nil, err
	}
	logx.Infow("[引擎] 节点执行完成",
		logx.Field("节点ID", nodeID),
		logx.Field("输出", output))
	return result, nil
}

// validateComponent 验证组件
func (e *WorkflowEngine) validateComponent(component components.Component, node *core.NodeDefinition) error {
	if validateErrors := component.Validate(); len(validateErrors) > 0 {
		return errors.New("组件 [" + node.ID + "] 验证失败: " + fmt.Sprintf("%+v", validateErrors))
	}
	return nil
}

// executeComponent 执行组件
func (e *WorkflowEngine) executeComponent(ctx *core.ExecutionContext, err error, component components.Component, input any) (*core.Result, error) {
	if err != nil {
		return nil, err
	}
	return component.Execute(ctx, input)
}

// updateNodeContext 更新节点上下文
func (e *WorkflowEngine) updateNodeContext(ctx *core.ExecutionContext, err error, nodeID string, output any) error {
	if err != nil {
		return err
	}
	if output != nil {
		logx.Debugf("[工作流] 更新节点上下文: %s, %+v", nodeID, output)
		ctx.SetVariable(nodeID+".output", output)
	}
	return nil
}

// createNodeResult 创建节点结果
func (e *WorkflowEngine) createNodeResult(node *core.NodeDefinition, input, output any, result *core.Result, startTime time.Time) *core.NodeResult {
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
func (e *WorkflowEngine) handleNodeError(ctx *core.ExecutionContext, node *core.NodeDefinition, err error, startTime time.Time) error {
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
		// fmt.Printf("更新工作流状态: %s, %+v\n", result.NodeID, result)
		ctx.SetNodeResult(result.NodeID, result)
	}

	// 计算整体进度
	completed := len(ctx.State.Result)
	total := int(ctx.TotalNodes)
	ctx.State.Progress = float64(completed) / float64(total)

	if err != nil {
		logx.Errorw("[引擎] 执行工作流失败",
			logx.Field("工作流ID", ctx.WorkspaceId),
			logx.Field("错误", err))
		return err
	}
	logx.Infow("[引擎] 工作流执行完成",
		logx.Field("工作流ID", ctx.WorkspaceId),
		logx.Field("输出", ctx.State.Result))
	return nil
}

// GetNodeResult 获取指定工作流节点的执行结果
func (e *WorkflowEngine) GetNodeResult(workflowID, serialID, nodeID string) (*core.NodeResult, bool) {
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()

	if !ok {
		logx.Debugf("[工作流] 获取节点结果, 工作流未找到: %s", workflowID)
		return nil, false
	}

	// 使用sync.Map获取执行上下文
	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		logx.Debugf("[工作流] 获取节点结果, 执行上下文未找到: %s", serialID)
		return nil, false
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		logx.Debugf("[工作流] 获取节点结果, 执行上下文类型错误: %s", serialID)
		return nil, false
	}

	return execCtx.GetNodeResult(nodeID)
}

// buildExecutionPlan 构建执行计划
func (e *WorkflowEngine) buildExecutionPlan(ctx context.Context, def *core.WorkflowDef) (*ExecutionPlan, error) {
	// 构建依赖图
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// 初始化入度
	for nodeID, def := range def.Nodes {
		inDegree[nodeID] = 0
		graph[nodeID] = []string{}
		// 迭代组件初始化
		if def.Type == "iteration" {
			logx.Debugf("[工作流] 迭代组件初始化: %s", def.ID)
			err := e.RegisterWorkflow(ctx, def.SubWorkflow)
			if err != nil {
				return nil, errors.New("迭代组件初始化失败: " + err.Error())
			}
		}
	}

	// 构建图结构
	for _, conn := range def.Connections {
		graph[conn.From] = append(graph[conn.From], conn.To)
		inDegree[conn.To]++
	}

	// 拓扑排序构建执行阶段
	var phases []ExecutionPhase
	for len(inDegree) > 0 {
		var phaseNodes []*WorkflowNode
		for nodeID, degree := range inDegree {
			if degree == 0 {
				node := def.Nodes[nodeID]
				phaseNodes = append(phaseNodes, &WorkflowNode{
					NodeDefinition: node,
				})
				delete(inDegree, nodeID)
			}
		}

		if len(phaseNodes) == 0 {
			return nil, errors.New("工作流中存在循环依赖")
		}

		phases = append(phases, ExecutionPhase{Nodes: phaseNodes})

		// 更新入度
		for _, node := range phaseNodes {
			for _, next := range graph[node.ID] {
				inDegree[next]--
			}
		}
	}
	return &ExecutionPlan{Phases: phases}, nil
}

// buildConditionRouter 构建条件路由
func (e *WorkflowEngine) buildConditionRouter(def *core.WorkflowDef) (map[string][]string, error) {
	router := make(map[string][]string)
	for _, conn := range def.Connections {
		router[conn.To] = append(router[conn.To], conn.From+"_"+conn.Condition)
	}
	return router, nil
}

// 启动清理协程
func (e *WorkflowEngine) startCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(e.config.CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.cleanupExecContexts()
				e.adjustPoolSize()
			case workflowID := <-e.cleanup:
				e.mu.Lock()
				delete(e.executorPool, workflowID)
				e.mu.Unlock()
			case <-e.shutdownCh:
				return
			}
		}
	}()
}

// 调整池大小
func (e *WorkflowEngine) adjustPoolSize() {
	// 获取当前运行的任务数
	running := e.pool.Running()
	cap := e.pool.Cap()

	// 如果使用率超过80%，且未达到最大值，则增加容量
	if float64(running)/float64(cap) > 0.8 && cap < e.config.MaxPoolSize {
		newCap := int(float64(cap) * 1.5)
		if newCap > e.config.MaxPoolSize {
			newCap = e.config.MaxPoolSize
		}
		e.pool.Tune(newCap)
		e.poolSize = newCap
	}

	// 如果使用率低于30%，且池大小超过初始值，则减少容量
	if running > 0 && float64(running)/float64(cap) < 0.3 && cap > e.config.InitialPoolSize {
		newCap := int(float64(cap) * 0.7)
		if newCap < e.config.InitialPoolSize {
			newCap = e.config.InitialPoolSize
		}
		e.pool.Tune(newCap)
		e.poolSize = newCap
	}
}

// cleanupExecContexts 清理执行上下文
func (e *WorkflowEngine) cleanupExecContexts() {
	now := time.Now()

	e.mu.RLock()
	executors := make([]*Executor, 0, len(e.executorPool))
	for _, exec := range e.executorPool {
		executors = append(executors, exec)
	}
	e.mu.RUnlock()

	for _, executor := range executors {
		// 清理每个执行器的上下文
		executor.execContexts.Range(func(key, value any) bool {
			id := key.(string)
			ctx, ok := value.(*core.ExecutionContext)
			if !ok {
				executor.execContexts.Delete(key)
				return true
			}

			// 检查是否已完成或过期
			if ctx.State.Status == core.StatusCompleted ||
				ctx.State.Status == core.StatusFailed ||
				now.After(ctx.Expiration) {
				executor.execContexts.Delete(id)
				core.ReleaseExecutionContext(ctx) // 释放上下文
			}
			return true
		})
	}
}

// Cleanup 释放资源
func (e *WorkflowEngine) Cleanup() {
	close(e.shutdownCh)
	e.pool.Release()

	// 关闭所有执行器
	e.mu.Lock()
	for id, executor := range e.executorPool {
		close(executor.shutdownCh)
		delete(e.executorPool, id)
		logx.Debugf("[工作流] 关闭执行器: %s", executor.definition.ID)
	}
	e.mu.Unlock()
}

// GetPoolStats 获取池统计信息
func (e *WorkflowEngine) GetPoolStats() map[string]any {
	return map[string]any{
		"capacity":    e.pool.Cap(),
		"running":     e.pool.Running(),
		"free":        e.pool.Free(),
		"waiting":     e.pool.Waiting(),
		"poolSize":    e.poolSize,
		"maxPoolSize": e.config.MaxPoolSize,
		"initialSize": e.config.InitialPoolSize,
	}
}

// GetExecutionContext 获取执行上下文
func (e *WorkflowEngine) GetExecutionContext(workflowID, serialID string) (*core.ExecutionContext, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	executor, ok := e.executorPool[workflowID]
	if !ok {
		return nil, errors.New("工作流未找到: " + workflowID)
	}

	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return nil, errors.New("执行上下文未找到: " + workflowID + ", " + serialID)
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return nil, errors.New("执行上下文类型错误: " + workflowID + ", " + serialID)
	}

	return execCtx, nil
}

func (e *WorkflowEngine) ClearExecutionContext(workflowID, serialID string) error {
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()
	if !ok {
		return errors.New("工作流未找到: " + workflowID)
	}

	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return errors.New("执行上下文未找到: " + workflowID + ", " + serialID)
	}

	execCtx, ok := val.(*core.ExecutionContext)
	executor.execContexts.Delete(serialID)
	core.ReleaseExecutionContext(execCtx)

	return nil
}

// ListWorkflows 列出所有工作流
func (e *WorkflowEngine) ListWorkflows() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	workflowIDs := make([]string, 0, len(e.executorPool))
	for id := range e.executorPool {
		workflowIDs = append(workflowIDs, id)
	}

	return workflowIDs
}

// GetWorkflowStatus 获取工作流状态
func (e *WorkflowEngine) GetWorkflowStatus(workflowID string) (core.WorkflowStatus, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	executor, ok := e.executorPool[workflowID]
	if !ok {
		return "", false
	}

	return executor.status, true
}

// PauseWorkflow 暂停工作流执行
func (e *WorkflowEngine) PauseWorkflow(ctx context.Context, workflowID string, serialID string) error {
	// 获取工作流执行器
	e.mu.RLock()
	executor, ok := e.executorPool[workflowID]
	e.mu.RUnlock()

	if !ok {
		return errors.New("工作流未找到: " + workflowID)
	}

	// 检查工作流状态
	if executor.status == core.WorkflowStatusShutdown {
		return errors.New("工作流正在关闭: " + workflowID)
	}

	// 使用sync.Map获取执行上下文
	val, ok := executor.execContexts.Load(serialID)
	if !ok {
		return errors.New("工作流执行实例未找到: " + workflowID + ", " + serialID)
	}

	execCtx, ok := val.(*core.ExecutionContext)
	if !ok {
		return errors.New("工作流执行上下文类型错误: " + workflowID + ", " + serialID)
	}

	// 检查执行状态
	if execCtx.State.Status != core.StatusRunning {
		return errors.New("工作流执行状态不是运行中: " + workflowID + ", " + serialID + ", 当前状态: " + string(execCtx.State.Status))
	}

	// 更新执行状态为暂停
	execCtx.State.Status = core.StatusPaused

	// 获取并调用取消函数
	cancelVal, ok := execCtx.GetVariable("cancel")
	if !ok {
		return errors.New("无法获取取消函数: " + workflowID + ", " + serialID)
	}

	cancel, ok := cancelVal.(context.CancelFunc)
	if !ok {
		return errors.New("取消函数类型错误: " + workflowID + ", " + serialID)
	}

	// 调用取消函数
	cancel()
	return nil
}
