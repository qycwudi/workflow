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
	// canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	// if err != nil {
	// 	logx.Errorw("[画布] 获取工作流定义失败",
	// 		logx.Field("工作空间ID", req.Id),
	// 		logx.Field("错误", err))
	// 	return nil, errors.New(int(logic.SystemOrmError), "获取工作流定义失败")
	// }

	// flowgram 测试 default 工作流
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		logx.Errorw("[画布] 获取工作流定义失败",
			logx.Field("工作空间ID", req.Id),
			logx.Field("错误", err))
		return nil, errors.New(int(logic.SystemOrmError), "获取工作流定义失败")
	}
	// 注册任务流
	err = workflow.Register(l.ctx, canvas.WorkspaceId, canvas.Draft)
	if err != nil {
		logx.Errorw("[画布] 注册工作流失败",
			logx.Field("工作空间ID", req.Id),
			logx.Field("错误", err))
		return nil, errors.New(int(logic.SystemError), "注册工作流失败")
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
		logx.Errorw("Failed to save the running record", logx.Field("error", err))
	}

	resp = &types.CanvasRunResponse{
		Id:       req.Id,
		SerialId: traceId,
	}
	go run(context.Background(), l.svcCtx, traceId, canvas.WorkspaceId, data)
	return resp, nil
}

func run(ctx context.Context, svcCtx *svc.ServiceContext, traceId, workspaceId string, data map[string]any) {
	var executeErr error
	// 运行文件
	_, result, err := workflow.Run(ctx, traceId, workspaceId, data)
	if err != nil {
		executeErr = fmt.Errorf("[画布] 执行工作流失败 [工作空间ID:%s] [序列ID:%s] [错误:%v]", workspaceId, traceId, err)
		logx.Errorw("[画布] 执行工作流失败",
			logx.Field("工作空间ID", workspaceId),
			logx.Field("序列ID", traceId),
			logx.Field("错误", err))
	} else {
		logx.Infow("Run the task flow successfully", logx.Field("result", result))
	}
	record, _ := svcCtx.SpaceRecordModel.FindOneBySerialNumber(ctx, traceId)
	other, _ := sonic.Marshal(result.Output)
	record.Other = string(other)
	status := enum.RecordStatusSuccess
	if result.Error != "" {
		status = enum.RecordStatusFail
		other, _ = sonic.Marshal(map[string]string{
			"error": result.Error,
		})
		record.Other = string(other)
	}
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
		logx.Errorw("更新运行记录失败", logx.Field("error", err))
	}
}
