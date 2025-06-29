package canvas

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/x/errors"

	"workflow/internal/enum"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
	"workflow/internal/workflow"
	pkgutils "workflow/pkg/utils"
)

type CanvasRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunLogic {
	return &CanvasRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunLogic) CanvasRun(req *types.CanvasRunRequest) (resp *types.CanvasRunResponse, err error) {
	traceId := "trace-" + trace.TraceIDFromContext(l.ctx)
	startTime := time.Now()
	
	// 记录工作流执行开始的审计日志
	pkgutils.LogBusinessOperation(pkgutils.OpExecute, req.Id, "", logx.Field("trace_id", traceId), logx.Field("operation", "canvas_run_start"))

	// flowgram 测试 default 工作流
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		pkgutils.LogModuleError(pkgutils.ModuleCanvas, "canvas_lookup", err, logx.Field("workspace_id", req.Id))
		return nil, errors.New(int(logic.SystemOrmError), "查询画布失败")
	}
	// 注册任务流
	err = workflow.Register(l.ctx, canvas.WorkspaceId, canvas.Draft)
	if err != nil {
		pkgutils.LogModuleError(pkgutils.ModuleCanvas, "workflow_registration", err, logx.Field("workspace_id", req.Id))
		return nil, errors.New(int(logic.SystemError), "工作流注册失败")
	}

	// 读取 参数
	data := req.Params

	spaceRecord := model.SpaceRecord{
		WorkspaceId:  canvas.WorkspaceId,
		SerialNumber: traceId,
		RunTime:      startTime,
		RecordName:   utils.FormatDate(startTime),
		Duration:     0,
		Other:        "{}",
		Status:       enum.RecordStatusRunning,
	}
	// 记录
	_, err = l.svcCtx.SpaceRecordModel.Insert(l.ctx, &spaceRecord)
	if err != nil {
		pkgutils.LogModuleError(pkgutils.ModuleCanvas, "space_record_save", err, logx.Field("trace_id", traceId))
	}

	resp = &types.CanvasRunResponse{
		Id:       req.Id,
		SerialId: traceId,
	}
	// 记录工作流执行启动的审计日志
	pkgutils.LogBusinessOperation(pkgutils.OpExecute, canvas.WorkspaceId, "", 
		logx.Field("trace_id", traceId),
		logx.Field("operation", "canvas_run_initiated"),
		logx.Field("params_count", len(data)))
	
	go run(context.Background(), l.svcCtx, traceId, canvas.WorkspaceId, data)
	return resp, nil
}

func run(ctx context.Context, svcCtx *svc.ServiceContext, traceId, workspaceId string, data map[string]any) {
	var executeErr error
	// 运行文件
	_, result, err := workflow.Run(ctx, traceId, workspaceId, data)
	if err != nil {
		executeErr = fmt.Errorf("[canvas] execute workflow failed [workspace id:%s] [trace id:%s] [error:%v]", workspaceId, traceId, err)
		pkgutils.LogModuleError(pkgutils.ModuleCanvas, "workflow_execution", err, 
			logx.Field("workspace_id", workspaceId),
			logx.Field("trace_id", traceId))
		// 记录工作流执行失败的审计日志
		pkgutils.LogBusinessOperation(pkgutils.OpExecute, workspaceId, "", 
			logx.Field("trace_id", traceId),
			logx.Field("operation", "canvas_run_failed"),
			logx.Field("error", err.Error()))
	} else {
		// 记录工作流执行成功的审计日志
		pkgutils.LogBusinessOperation(pkgutils.OpExecute, workspaceId, "", 
			logx.Field("trace_id", traceId),
			logx.Field("operation", "canvas_run_completed"),
			logx.Field("result_size", len(fmt.Sprintf("%v", result))))
	}
	record, _ := svcCtx.SpaceRecordModel.FindOneBySerialNumber(ctx, traceId)
	other, _ := sonic.Marshal(result)
	record.Other = string(other)
	status := enum.RecordStatusSuccess

	if executeErr != nil {
		record.Status = enum.RecordStatusFail
		other, _ = sonic.Marshal(map[string]string{
			"error": executeErr.Error(),
		})
		record.Other = string(other)
	}
	record.Status = status

	record.Duration = time.Since(record.RunTime).Milliseconds()

	err = svcCtx.SpaceRecordModel.Update(ctx, record)
	if err != nil {
		pkgutils.LogModuleError(pkgutils.ModuleDatabase, "space_record_update", err, 
			logx.Field("trace_id", traceId))
	}
}
