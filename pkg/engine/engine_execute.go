package engine

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/components"
	"workflow/pkg/constants"
	"workflow/pkg/core"
	"workflow/pkg/utils"
)

// ExecuteWorkflow 执行工作流
func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID string, traceID string, params map[string]any, extra core.ContextExtra) (*core.ExecutionContextEnhanced, error) {
	// 参数校验
	if params == nil {
		return nil, &EngineExecuteError{Message: "params is nil"}
	}

	utils.LogDebugInfo(utils.ModuleWorkflow, "workflow_execution", map[string]any{"params_count": len(params)})

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
	ExecutionContextEnhanced := core.NewExecutionContextEnhanced(ctx, workflowID, traceID, params, extra, runnable.defaultTTL, runnable.nodesNum)

	// 存储执行上下文
	runnable.execContexts.Store(traceID, ExecutionContextEnhanced)

	// 记录工作流开始追踪
	startTime := time.Now()
	if ExecutionContextEnhanced.IsTrace && trace != nil {
		workflowTrace := &TraceRecore{
			WorkspaceId: ExecutionContextEnhanced.WorkspaceId,
			TraceId:     traceID,
			NodeId:      "workflow",
			NodeName:    "工作流执行",
			NodeType:    "workflow",
			Status:      "running",
			StartTime:   startTime,
			Input:       params,
			Step:        0,
		}
		trace.CreateTrace(ctx, workflowTrace)
	}

	// 执行工作流
	err := e.executeWorkflowPhases(ExecutionContextEnhanced, runnable)

	// 确保调用取消函数
	ExecutionContextEnhanced.Cancel()

	// 记录工作流完成追踪
	if ExecutionContextEnhanced.IsTrace && trace != nil {
		status := "completed"
		errorMsg := ""
		var output any

		if err != nil {
			status = "failed"
			errorMsg = err.Error()
		} else {
			// 获取工作流输出
			if endOutput, ok := ExecutionContextEnhanced.GetVariable(constants.EndParameters); ok {
				output = endOutput
			}
		}

		workflowTrace := &TraceRecore{
			WorkspaceId: ExecutionContextEnhanced.WorkspaceId,
			TraceId:     traceID,
			NodeId:      "workflow",
			NodeName:    "工作流执行",
			NodeType:    "workflow",
			Status:      status,
			Output:      output,
			ErrorMsg:    errorMsg,
			ElapsedTime: time.Since(startTime).Milliseconds(),
		}
		trace.UpdateTrace(ctx, workflowTrace)
	}

	// 处理执行结果
	if err != nil {
		return nil, err
	}
	return ExecutionContextEnhanced, nil
}

func (e *WorkflowEngine) ExecuteSingleWorkflow(ctx context.Context, workflowID, serialID string, nodeId string, params map[string]any) (*core.NodeResult, error) {
	startTime := time.Now()

	// 获取工作流执行器
	e.mu.RLock()
	runnable, ok := e.runnablePool[workflowID]
	e.mu.RUnlock()

	if !ok {
		return nil, &EngineExecuteError{Message: "workflow not found: " + workflowID}
	}

	// 查找目标节点
	var node *core.Nodes
	for _, n := range runnable.definition.Nodes {
		if n.ID == nodeId {
			node = &n
			break
		}
		if n.Type == constants.ComponentLoop {
			for _, block := range n.Blocks {
				if block.ID == nodeId {
					node = &block
					break
				}
			}
		}
	}

	if node == nil {
		return nil, &EngineExecuteError{Message: "node not found: " + nodeId}
	}

	// 创建节点结果容器
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
		IsSingleNode:      true, // 标记为单节点执行
	}
	execCtx := core.NewExecutionContextEnhanced(ctx, workflowID, serialID, params, extra, runnable.defaultTTL, runnable.nodesNum)

	// 创建执行管理器
	executionManager := NewExecutionManager(1, e)
	defer executionManager.Close()

	// 创建工作流节点
	workflowNode := &WorkflowNode{
		Nodes: node,
	}

	// 执行节点
	err := executionManager.executeNodeSimplified(execCtx, runnable, workflowNode)
	if err != nil {
		utils.LogModuleError(utils.ModuleWorkflow, "single_node_execution", err, logx.Field("node_id", nodeId))
		nodeResult.Error = err.Error()
		return nodeResult, &EngineExecuteError{Message: "failed to execute node: " + err.Error()}
	}

	// 获取节点输出
	var output any
	var outputExists bool

	if node.Type == constants.ComponentEnd {
		output, outputExists = execCtx.GetVariable(constants.EndParameters)
	} else {
		output, outputExists = execCtx.GetVariable(nodeId + ".output")
	}

	if !outputExists {
		return nil, &EngineExecuteError{Message: "node output not found"}
	}

	// 从执行结果中保持原有路由信息
	// 路由信息在组件执行过程中已通过SetRoute设置到ExecutionContext中

	nodeResult.Output = output
	nodeResult.Duration = time.Since(startTime).Milliseconds()
	return nodeResult, nil
}

// executeWorkflowPhases 执行工作流阶段
func (e *WorkflowEngine) executeWorkflowPhases(execCtx *core.ExecutionContextEnhanced, runnable *Runnable) error {
	// 创建执行管理器
	executionManager := NewExecutionManager(1000, e)
	defer executionManager.Close()

	// 记录工作流执行开始
	utils.LogWorkflowExecution(execCtx.WorkspaceId, execCtx.WorkspaceId, execCtx.TraceId, "", "started", 0, int(runnable.nodesNum))
	for _, phase := range runnable.executionPlan.Phases {
		if err := executionManager.ExecutePhaseSimplified(execCtx, runnable, phase); err != nil {
			utils.LogModuleError(utils.ModuleWorkflow, "phase_execution", err, logx.Field("phase_index", len(runnable.executionPlan.Phases)))
			return err
		}
	}
	// 记录工作流执行完成
	startTime := time.Now()
	utils.LogWorkflowExecution(execCtx.WorkspaceId, execCtx.WorkspaceId, execCtx.TraceId, "", "completed", time.Since(startTime).Milliseconds(), int(runnable.nodesNum))
	return nil
}
