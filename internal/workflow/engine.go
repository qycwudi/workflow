package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/pkg/core"
	"workflow/pkg/engine"
	"workflow/pkg/metrics"
)

var eg *engine.WorkflowEngine

func InitEngine(ctx *svc.ServiceContext) {
	// 创建工作流引擎
	eg = engine.NewWorkflowEngine()
	defer eg.Cleanup() // 确保在程序结束时释放资源
	metrics.NewCollector()
}

func Register(ctx context.Context, dsl string) error {
	workflowDefinition := core.WorkflowDef{}
	err := json.Unmarshal([]byte(dsl), &workflowDefinition)
	if err != nil {
		logx.Errorw("解析工作流文件失败", logx.Field("error", err))
		return err
	}

	definitionBytes, _ := json.MarshalIndent(workflowDefinition, "", "  ")
	logx.Infow("工作流定义", logx.Field("definition", string(definitionBytes)))

	err = eg.RegisterWorkflow(ctx, &workflowDefinition)
	if err != nil {
		logx.Errorw("注册工作流失败", logx.Field("error", err))
		return err
	}
	return nil
}

func Run(ctx context.Context, workspaceId string, data map[string]any) (string, core.NodeResult, error) {
	serialId := uuid.New().String()
	// 执行工作流
	if err := eg.ExecuteWorkflow(ctx, workspaceId, serialId, data); err != nil {
		log.Fatalf("工作流执行失败: %v", err)
		return serialId, core.NodeResult{}, err
	}

	endResult, ok := eg.GetNodeResult(workspaceId, serialId, "end-node-1")
	if !ok {
		logx.Errorw("未找到 end-node-1 节点的执行结果")
		return serialId, core.NodeResult{}, errors.New("未找到 end-node-1 节点的执行结果")
	}
	return serialId, *endResult, nil
}
