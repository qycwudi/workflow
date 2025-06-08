package engine

import (
	"context"

	"github.com/rotisserie/eris"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/components"
	"workflow/pkg/core"
)

// RegisterWorkflow 注册工作流
func (e *WorkflowEngine) RegisterWorkflow(ctx context.Context, id string, def *core.WorkflowDef) error {

	// 构建执行计划
	plan, err := e.buildExecutionPlan(ctx, def)
	if err != nil {
		return eris.New("failed to build execution plan: " + err.Error())
	}

	// 构建条件路由
	condition, err := e.buildConditionRouter(def)
	if err != nil {
		return eris.New("failed to build condition mapping: " + err.Error())
	}

	// 创建新的执行器
	executor := e.createExecutor(id, def, plan, condition)

	// 注册执行器
	e.registerExecutor(id, executor)
	return nil
}

// buildExecutionPlan 构建执行计划
func (e *WorkflowEngine) buildExecutionPlan(ctx context.Context, def *core.WorkflowDef) (*ExecutionPlan, error) {
	// 构建依赖图
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	nodeMap := make(map[string]*core.Nodes)
	// 初始化入度
	for _, node := range def.Nodes {
		nodeMap[node.ID] = &node
		inDegree[node.ID] = 0
		graph[node.ID] = []string{}
		// 迭代组件初始化
		if node.Type == components.Loop {
			logx.Debugf("[工作流] 迭代组件初始化: %s", node.ID)
			err := e.RegisterWorkflow(ctx, node.ID, &core.WorkflowDef{
				Nodes: node.Blocks,
				Edges: node.Edges,
			})
			if err != nil {
				return nil, eris.New("迭代组件初始化失败: " + err.Error())
			}
		}
	}

	// 构建图结构
	for _, conn := range def.Edges {
		graph[conn.SourceNodeID] = append(graph[conn.SourceNodeID], conn.TargetNodeID)
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
			return nil, eris.New("workflow has circular dependency")
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
	for _, conn := range def.Edges {
		if conn.SourcePortID == "" {
			// 如果条件为空,则默认成功
			conn.SourcePortID = components.Success
		}
		router[conn.TargetNodeID] = append(router[conn.TargetNodeID], conn.SourceNodeID+"_"+conn.SourcePortID)
	}
	return router, nil
}

// registerExecutor 注册执行器
func (e *WorkflowEngine) registerExecutor(workflowID string, executor *Executor) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.executorPool[workflowID] = executor
}
