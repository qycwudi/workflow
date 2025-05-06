package workflow

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/pkg/core"
	"workflow/pkg/engine"
)

var eg *engine.WorkflowEngine

func InitEngine(ctx *svc.ServiceContext) {
	// 创建工作流引擎
	eg = engine.NewWorkflowEngine()
}

func Register(ctx context.Context, dsl string) error {
	var workflowDefinition core.WorkflowDef
	err := json.Unmarshal([]byte(dsl), &workflowDefinition)
	if err != nil {
		logx.Errorw("解析工作流文件失败", logx.Field("error", err))
		return err
	}

	logx.Debugw("工作流DSL", logx.Field("dsl", dsl))
	definitionBytes, _ := json.MarshalIndent(workflowDefinition, "", "  ")
	logx.Infow("工作流定义", logx.Field("definition", string(definitionBytes)))

	err = eg.RegisterWorkflow(ctx, &workflowDefinition)
	if err != nil {
		logx.Errorw("注册工作流失败", logx.Field("error", err))
		return err
	}
	return nil
}

func Run(ctx context.Context, serialId, workspaceId string, data map[string]any) (string, core.NodeResult, error) {
	defer func() {
		logx.Infow("清除执行上下文", logx.Field("workspaceId", workspaceId), logx.Field("serialId", serialId))
		clearErr := eg.ClearExecutionContext(workspaceId, serialId)
		if clearErr != nil {
			logx.Errorw("清除执行上下文失败", logx.Field("error", clearErr.Error()))
		}
		logx.Infow("清除执行上下文成功", logx.Field("workspaceId", workspaceId), logx.Field("serialId", serialId))
	}()
	// 执行工作流
	if err := eg.ExecuteWorkflow(ctx, workspaceId, serialId, data); err != nil {
		logx.Errorw("工作流执行失败", logx.Field("error", err))
		return serialId, core.NodeResult{}, err
	}

	endResult, ok := eg.GetNodeResult(workspaceId, serialId, "end-node-1")
	if !ok {
		logx.Errorw("未找到 结束 节点的执行结果", logx.Field("workspaceId", workspaceId), logx.Field("serialId", serialId))
		return serialId, core.NodeResult{}, errors.New("未找到 结束 节点的执行结果:" + workspaceId + "," + serialId)
	}
	return serialId, *endResult, nil
}

func RunSingle(ctx context.Context, serialId, workspaceId string, nodeId string, data map[string]any) (string, core.NodeResult, error) {
	// 执行工作流
	result, err := eg.ExecuteSingleWorkflow(ctx, workspaceId, serialId, nodeId, data)
	if err != nil {
		logx.Errorw("工作流执行失败", logx.Field("error", err))
		return serialId, *result, err
	}

	return serialId, *result, nil
}

func Close() {
	eg.Cleanup()
}

func Stop(ctx context.Context, workspaceId string, serialId string) error {
	return nil
}
