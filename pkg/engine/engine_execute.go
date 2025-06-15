package engine

import (
	"context"
	"fmt"
	"maps"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/components"
	"workflow/pkg/core"
)

// ExecuteWorkflow 执行工作流
func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID string, traceID string, params map[string]any, extra core.ContextExtra) (*core.ExecutionContext, error) {
	// 参数校验
	if params == nil {
		return nil, &EngineExecuteError{Message: "params is nil"}
	}

	logx.Debugf("[Workflow] Execute parameters: %+v", params)

	// 获取工作流执行器
	e.mu.RLock()
	runnable, ok := e.runnablePool[workflowID]
	e.mu.RUnlock()

	if !ok {
		return nil, &EngineExecuteError{Message: "workflow not found: " + workflowID}
	}

	// 检查工作流状态
	if runnable.status == core.WorkflowStatusShutdown {
		return nil, &EngineExecuteError{Message: "workflow is shutting down: " + workflowID}
	}

	// 创建执行上下文
	executionContext := core.NewExecutionContext(ctx, workflowID, traceID, params, extra, runnable.defaultTTL, runnable.nodesNum)

	// 存储执行上下文
	runnable.execContexts.Store(traceID, executionContext)

	// 执行工作流
	err := e.executeWorkflowPhases(executionContext, runnable)

	// 确保调用取消函数
	executionContext.Cancel()
	// 处理执行结果
	if err != nil {
		return nil, err
	}
	return executionContext, nil
}

func (e *WorkflowEngine) ExecuteSingleWorkflow(ctx context.Context, workflowID, serialID string, nodeId string, params map[string]any) (*core.NodeResult, error) {
	startTime := time.Now()
	runnable, ok := e.runnablePool[workflowID]
	if !ok {
		return nil, &EngineExecuteError{Message: "workflow not found: " + workflowID}
	}

	var node *core.Nodes
	// 遍历获取节点
	for _, n := range runnable.definition.Nodes {
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
	extra := core.ContextExtra{
		IsSub:             false,
		ParentWorkspaceId: "",
		Index:             0,
		NodeNum:           runnable.nodesNum,
	}
	execCtx := core.NewExecutionContext(ctx, workflowID, serialID, params, extra, runnable.defaultTTL, runnable.nodesNum)

	// 执行节点
	component, err := components.ComponentFactory(e, node.Type, node.Data)
	if err != nil {
		logx.Errorw("failed to create component", logx.Field("error", err.Error()))
		nodeResult.Error = err.Error()
		return nodeResult, &EngineExecuteError{Message: "failed to create component: " + err.Error()}
	}
	err = e.executeNode(execCtx, &WorkflowNode{
		Component: component,
		Nodes:     node,
	})
	if err != nil {
		logx.Errorw("failed to execute node", logx.Field("error", err.Error()))
		nodeResult.Error = err.Error()
		return nodeResult, &EngineExecuteError{Message: "failed to execute node: " + err.Error()}
	}
	// 获取结束节点输出
	endOutput, ok := execCtx.GetVariable(core.EndParameters)
	if !ok {
		return nil, &EngineExecuteError{Message: "end node output not found"}
	}
	nodeResult.Output = endOutput
	nodeResult.Duration = time.Since(startTime).Milliseconds()
	return nodeResult, nil
}

// executeWorkflowPhases 执行工作流阶段
func (e *WorkflowEngine) executeWorkflowPhases(execCtx *core.ExecutionContext, runnable *Runnable) error {
	// 总执行计划
	logx.Debugf("[Workflow] Total execution plan: %d", len(runnable.executionPlan.Phases))
	for _, phase := range runnable.executionPlan.Phases {
		if err := e.executePhase(execCtx, runnable, phase); err != nil {
			logx.Errorw("[Workflow] Execute phase failed", logx.Field("error", err.Error()))
			return err
		}
	}
	return nil
}

// executePhase 执行单个阶段
func (e *WorkflowEngine) executePhase(execCtx *core.ExecutionContext, runnable *Runnable, phase ExecutionPhase) error {

	logx.Debugf("[Workflow] execute phase node count: %d", len(phase.Nodes))
	// 获取单例协程池
	var execWg sync.WaitGroup
	pool := GetGlobalPool()

	// 创建错误收集器
	var errMu sync.Mutex
	var errorStack []error

	for _, node := range phase.Nodes {
		execWg.Add(1)
		err := pool.Submit(func() {
			defer execWg.Done()
			err := e.handleNodeExecution(execCtx, runnable, node)
			if err != nil {
				logx.Errorf("[Workflow] Execute node failed, nodeId:%s, nodeName:%s, nodeType:%s, traceId:%s, error %s",
					node.ID, node.Data.Title, node.Type, execCtx.TraceId, err.Error())
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
		return fmt.Errorf("[workflow_execute] phase execute failed: %v", errorStack)
	}
	return nil
}

// handleNodeExecution 处理节点执行
func (e *WorkflowEngine) handleNodeExecution(execCtx *core.ExecutionContext, runnable *Runnable, node *WorkflowNode) error {
	// 开始类型的组件跳过路由检查
	if node.Type != components.Start {
		// 检查路由
		ok, err := e.checkNodeRoute(execCtx, runnable, node)
		if err != nil {
			return err
		}
		if !ok {
			logx.Debugf("[workflow_execute] node route skipped: %s", node.ID)
			return nil
		}
	}

	// 创建并执行组件
	component, err := components.ComponentFactory(e, node.Type, node.Data)
	if err != nil {
		return &EngineExecuteError{Message: "failed to create component [" + node.ID + "]: " + err.Error()}
	}

	node.Component = component

	err = e.executeNode(execCtx, node)
	if err != nil {
		return err
	}
	return nil
}

// checkNodeRoute 检查节点路由
func (e *WorkflowEngine) checkNodeRoute(execCtx *core.ExecutionContext, executor *Runnable, node *WorkflowNode) (bool, error) {
	route, ok := executor.conditionRouter[node.ID]
	if !ok {
		logx.Debugf("[Workflow] Node route not found, skip[loop single node]: %s", node.ID)
		return true, nil
	}

	if !execCtx.CheckRoute(route) {
		logx.Debugf("[workflow_execute] skip node route:%+v", route)
		return false, e.handleSkippedNode(execCtx, node)
	}

	return true, nil
}

// handleSkippedNode 处理跳过的节点
func (e *WorkflowEngine) handleSkippedNode(execCtx *core.ExecutionContext, node *WorkflowNode) error {
	logx.Debugf("[workflow_execute] %s node necessary route not found, fill empty parameters", node.ID)
	// 补齐默认零值
	output := make(map[string]any)
	output, err := core.ProcessNodeOutput(output, node.Data.NodeDataOutputs)
	if err != nil {
		return err
	}
	execCtx.SetVariable(node.ID+".output", output)
	return nil
}

// executeNode 执行节点
func (e *WorkflowEngine) executeNode(execCtx *core.ExecutionContext, node *WorkflowNode) error {
	logx.Debugf("[Workflow] execute node: %s", node.ID)
	startTime := time.Now()
	var err error
	// 创建 trace 记录
	tid, err := e.createTrace(execCtx, node, startTime)
	if err != nil {
		return err
	}
	// 验证组件
	err = e.validateComponent(node.Component, node.ID, err)

	// 准备输入数据
	input, err := e.prepareNodeInput(execCtx, node, err)

	// 执行组件
	result, err := e.executeComponent(execCtx, node.Component, input, err)

	// 处理输出数据
	output, err := e.processNodeOutput(input, result, node, err)
	// 释放组件
	node.Component.Clear()

	// 更新 trace 记录
	e.updateTrace(execCtx, tid, node, startTime, input, output, err)

	// 更新结果
	if node.Type == components.End {
		// 设置结束节点输出
		execCtx.SetVariable(core.EndParameters, output)
	} else {
		// 设置标准节点输出
		execCtx.SetVariable(node.ID+".output", output)
	}
	// 更新路由
	execCtx.SetRoute(node.ID, result.Route)
	logx.Debugf("[Workflow] execute node %s,input:%+v,output:%+v,error:%+v", node.ID, input, output, err)
	return nil
}

// prepareNodeInput 准备节点输入数据
func (e *WorkflowEngine) prepareNodeInput(execCtx *core.ExecutionContext, node *WorkflowNode, err error) (map[string]any, error) {
	if err != nil {
		return nil, err
	}

	if node.Type == components.Start {
		zero, ok := execCtx.GetVariable(core.GenesisParameters + execCtx.WorkspaceId)
		if !ok {
			logx.Errorf("[workflow_execute] input parameter not found: %s", core.GenesisParameters+execCtx.WorkspaceId)
			return nil, &EngineExecuteError{Message: "input parameter not found"}
		}
		return core.ProcessNodeOutput(zero, node.Data.NodeDataOutputs)
	}
	if node.Type == components.End {
		// 结束节点 输入和输出相同,DSL 定义了输出, 所以需要将输出赋值给输入
		node.Data.NodeDataInputs.Properties = node.Data.NodeDataOutputs.Properties
		node.Data.NodeDataInputs.Required = node.Data.NodeDataOutputs.Required
		node.Data.NodeDataInputs.Type = node.Data.NodeDataOutputs.Type
	}

	// 初始化合并结果
	mergedInput := make(map[string]any)

	// 1. 获取自定义输入
	customInput, err := node.Component.AnalyzeInputs(execCtx)
	if err != nil {
		logx.Errorf("[workflow_execute] failed to analyze custom input: %s", err.Error())
		return nil, &EngineExecuteError{Message: "failed to analyze custom input: " + err.Error()}
	}

	// 2. 如果有自定义输入，先复制到合并结果中
	if customInput != nil {
		if customMap, ok := customInput.(map[string]any); ok {
			for k, v := range customMap {
				mergedInput[k] = v
			}
		} else {
			logx.Errorf("[workflow_execute] custom input must be map type")
			return nil, &EngineExecuteError{Message: "custom input must be map type"}
		}
	}

	// 3. 获取标准输入
	standardInput, err := core.ParseNodeInputs(execCtx, node.Data.NodeDataInputsValue, node.Data.NodeDataInputs)
	if err != nil {
		logx.Errorf("[workflow_execute] failed to parse standard input: %s", err.Error())
		return nil, &EngineExecuteError{Message: "failed to parse standard input: " + err.Error()}
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
func (e *WorkflowEngine) processNodeOutput(input map[string]any, result *core.Result, node *WorkflowNode, err error) (map[string]any, error) {
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
		logx.Errorf("[workflow_execute] failed to process output mapping: %s", err.Error())
		return nil, &EngineExecuteError{Message: "failed to process output mapping: " + err.Error()}
	}
	return output, nil
}

func (e *WorkflowEngine) createTrace(ctx *core.ExecutionContext, node *WorkflowNode, startTime time.Time) (int64, error) {
	if !ctx.IsTrace {
		return 0, nil
	}
	trace := &TraceRecore{
		WorkspaceId: ctx.WorkspaceId,
		TraceId:     ctx.TraceId,
		NodeId:      node.ID,
		NodeName:    node.Data.Title,
		NodeType:    node.Type,
		Input:       "{}",
		Logic:       "{}",
		Output:      "{}",
		StartTime:   startTime,
		Step:        ctx.Extra.Index,
		SubIndex:    ctx.Extra.Index,
	}
	tid, err := Trace.CreateTrace(ctx, trace)
	if err != nil {
		logx.Errorf("[workflow_execute] create trace failed [Error:%v]", err)
	}
	return tid, nil
}

func (e *WorkflowEngine) updateTrace(ctx *core.ExecutionContext, tid int64, node *WorkflowNode, startTime time.Time, input, output map[string]any, err error) {
	if !ctx.IsTrace {
		return
	}
	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	}
	Trace.UpdateTraceById(ctx, &TraceRecore{
		Id:          tid,
		NodeId:      node.ID,
		TraceId:     ctx.TraceId,
		Input:       input,
		Output:      output,
		Status:      string(core.StatusCompleted),
		ElapsedTime: time.Since(startTime).Milliseconds(),
		ErrorMsg:    errorMsg,
	})
}

// validateComponent 验证组件
func (e *WorkflowEngine) validateComponent(component components.Component, nodeID string, err error) error {
	if err != nil {
		return err
	}
	validateErrors := component.Validate()
	if len(validateErrors) > 0 {
		return &EngineExecuteError{Message: "component [" + nodeID + "] validation failed: " + fmt.Sprintf("%+v", validateErrors)}
	}
	return nil
}

// executeComponent 执行组件
func (e *WorkflowEngine) executeComponent(ctx *core.ExecutionContext, component components.Component, input any, err error) (*core.Result, error) {
	if err != nil {
		return nil, err
	}
	return component.Execute(ctx, input)
}
