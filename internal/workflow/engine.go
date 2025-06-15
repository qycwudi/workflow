package workflow

import (
	"context"
	"errors"

	"github.com/bytedance/sonic"
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

func InitOpenApiEngine(ctx *svc.ServiceContext) {
	apis, err := ctx.ApiModel.FindByOn(context.Background())
	if err != nil {
		logx.Errorf("[workflow] find api error: %s", err)
		return
	}
	count := 0
	for _, api := range apis {
		Register(context.Background(), api.ApiId, api.Dsl)
		count++
	}
	logx.Infof("[workflow] register api success [count:%d]", count)
}

func Register(ctx context.Context, id string, dsl string) error {
	var graph core.Graph
	err := sonic.Unmarshal([]byte(dsl), &graph)
	if err != nil {
		logx.Errorf("[workflow] compile workflow error: %s", err)
		return err
	}
	err = eg.Compile(ctx, id, &graph)
	if err != nil {
		logx.Errorf("[workflow] register workflow error: %s", err)
		return err
	}
	return nil
}

func Run(ctx context.Context, traceId, workspaceId string, data map[string]any) (string, map[string]any, error) {
	logx.Infof("[workflow] start execute workflow [workspaceId:%s] [traceId:%s]", workspaceId, traceId)
	defer func() {
		logx.Infof("[workflow] clear execution context [workspaceId:%s] [traceId:%s]", workspaceId, traceId)
		clearErr := eg.ClearExecutionContext(workspaceId, traceId)
		if clearErr != nil {
			logx.Errorf("[workflow] clear execution context error: %s", clearErr.Error())
		}
	}()
	// 执行工作流
	execCtx, err := eg.ExecuteWorkflow(ctx, workspaceId, traceId, data, core.ContextExtra{
		IsSub:             false,
		ParentWorkspaceId: "",
		Index:             -1,
		NodeNum:           0,
	})
	if err != nil {
		logx.Errorf("[workflow] execute workflow error: %s", err)
		return traceId, nil, err
	}

	endResult, ok := execCtx.GetVariable(core.EndParameters)
	if !ok {
		logx.Errorf("[workflow] end node result not found [workspaceId:%s] [traceId:%s]", workspaceId, traceId)
		return traceId, nil, errors.New("end node result not found:" + workspaceId + "," + traceId)
	}
	return traceId, endResult, nil
}

func RunSingle(ctx context.Context, serialId, workspaceId string, nodeId string, data map[string]any) (string, core.NodeResult, error) {
	// 执行工作流
	result, err := eg.ExecuteSingleWorkflow(ctx, workspaceId, serialId, nodeId, data)
	if err != nil {
		logx.Errorf("[workflow] execute workflow error: %s", err)
		return serialId, *result, err
	}

	return serialId, *result, nil
}

func Stop(ctx context.Context, workspaceId string, serialId string) error {
	return nil
}
