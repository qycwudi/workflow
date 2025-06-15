package engine

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/components"
	"workflow/pkg/core"
)

// Compile 注册工作流
func (e *WorkflowEngine) Compile(ctx context.Context, id string, graph *core.Graph) error {

	// 构建执行计划
	plan, err := e.buildExecutionPlan(ctx, graph)
	if err != nil {
		return &EngineCompileError{
			Message: "failed to build execution plan: " + err.Error(),
		}
	}

	// 构建条件路由
	condition, err := e.buildConditionRouter(graph)
	if err != nil {
		return &EngineCompileError{
			Message: "failed to build condition mapping: " + err.Error(),
		}
	}

	// 创建新的执行器
	runnable := e.createRunnable(id, graph, plan, condition)

	// 注册执行器
	e.registerRunnable(id, runnable)
	return nil
}

// buildExecutionPlan 构建执行计划
func (e *WorkflowEngine) buildExecutionPlan(ctx context.Context, graph *core.Graph) (*ExecutionPlan, error) {
	// 构建依赖图
	plan := make(map[string][]string)
	inDegree := make(map[string]int)
	nodeMap := make(map[string]*core.Nodes)
	// 初始化入度
	for _, node := range graph.Nodes {
		nodeMap[node.ID] = &node
		inDegree[node.ID] = 0
		plan[node.ID] = []string{}
		// 迭代组件初始化
		if node.Type == components.Loop {
			logx.Debugf("[compile] loop component init: %s", node.ID)
			err := e.Compile(ctx, node.ID, &core.Graph{
				Nodes: node.Blocks,
				Edges: node.Edges,
			})
			if err != nil {
				return nil, &EngineCompileError{
					Message: "loop component init failed: " + err.Error(),
				}
			}
		}
	}

	// 构建图结构
	for _, conn := range graph.Edges {
		plan[conn.SourceNodeID] = append(plan[conn.SourceNodeID], conn.TargetNodeID)
		inDegree[conn.TargetNodeID]++
	}

	// 拓扑排序构建执行阶段
	var phases []ExecutionPhase
	for len(inDegree) > 0 {
		var phaseNodes []*WorkflowNode
		for nodeID, degree := range inDegree {
			if degree == 0 {
				node := nodeMap[nodeID]
				phaseNodes = append(phaseNodes, &WorkflowNode{
					Nodes: node,
				})
				delete(inDegree, nodeID)
			}
		}

		if len(phaseNodes) == 0 {
			return nil, &EngineCompileError{
				Message: "workflow has circular dependency",
			}
		}

		phases = append(phases, ExecutionPhase{Nodes: phaseNodes})

		// 更新入度
		for _, node := range phaseNodes {
			for _, next := range plan[node.ID] {
				inDegree[next]--
			}
		}
	}
	return &ExecutionPlan{Phases: phases}, nil
}

// buildConditionRouter 构建条件路由
func (e *WorkflowEngine) buildConditionRouter(graph *core.Graph) (map[string][]string, error) {
	router := make(map[string][]string)
	for _, conn := range graph.Edges {
		if conn.SourcePortID == "" {
			// if condition is empty, default to success
			conn.SourcePortID = components.Success
		}
		router[conn.TargetNodeID] = append(router[conn.TargetNodeID], conn.SourceNodeID+"_"+conn.SourcePortID)
	}
	return router, nil
}

// registerExecutor 注册执行器
func (e *WorkflowEngine) registerRunnable(workflowID string, runnable *Runnable) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.executorPool[workflowID] = runnable
}

// createExecutor 创建执行器
func (e *WorkflowEngine) createRunnable(id string, graph *core.Graph, plan *ExecutionPlan, condition map[string][]string) *Runnable {
	e.mu.Lock()
	defer e.mu.Unlock()
	runnable := &Runnable{
		definition:      graph,
		executionPlan:   plan,
		conditionRouter: condition,
		status:          core.WorkflowStatusActive,
		createdAt:       time.Now(),
		shutdownCh:      make(chan struct{}),
	}

	if existing, exists := e.executorPool[id]; exists {
		existing.status = core.WorkflowStatusDeploying
		// 复制已经存在的执行上下文,防止删除正在执行的上下文
		existing.execContexts.Range(func(key, value any) bool {
			runnable.execContexts.Store(key, value)
			return true
		})
	}
	return runnable
}
