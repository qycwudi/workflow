package canvas

import (
	"context"
	"time"

	"github.com/rulego/rulego/utils/json"
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

type CanvasRunSingleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunSingleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunSingleLogic {
	return &CanvasRunSingleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunSingleLogic) CanvasRunSingle(req *types.CanvasRunSingleRequest) (resp *types.CanvasRunSingleResponse, err error) {
	// 注册工作流
	traceId := "trace-" + trace.TraceIDFromContext(l.ctx)
	startTime := time.Now()
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		logx.Errorw("[画布] 获取工作流定义失败",
			logx.Field("工作空间ID", req.Id),
			logx.Field("错误", err))
		return nil, errors.New(int(logic.SystemOrmError), "获取工作流定义失败")
	}
	err = workflow.Register(l.ctx, canvas.Draft)
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
	record, err := l.svcCtx.SpaceRecordModel.Insert(l.ctx, &spaceRecord)
	if err != nil {
		logx.Errorw("Failed to save the running record", logx.Field("error", err))
	}
	spaceRecord.Id, _ = record.LastInsertId()
	// 查询工作流执行输入
	_, result, err := workflow.RunSingle(l.ctx, traceId, canvas.WorkspaceId, req.NodeId, data)
	if err != nil {
		logx.Errorw("[画布] 执行单节点失败",
			logx.Field("工作空间ID", req.Id),
			logx.Field("节点ID", req.NodeId),
			logx.Field("错误", err))
		spaceRecord.Status = enum.RecordStatusFail
	} else {
		spaceRecord.Status = enum.RecordStatusSuccess
	}
	other, _ := json.Marshal(result.Output)
	spaceRecord.Other = string(other)
	spaceRecord.Duration = time.Since(spaceRecord.RunTime).Milliseconds()
	err = l.svcCtx.SpaceRecordModel.Update(l.ctx, &spaceRecord)
	if err != nil {
		logx.Errorw("Failed to update the running record", logx.Field("error", err))
	}

	// 返回工作流执行结果 包括组件执行结果
	resp = &types.CanvasRunSingleResponse{
		Id:     req.Id,
		NodeId: req.NodeId,
		Result: types.Result{
			Input:     data,
			Output:    result.Output,
			NodeId:    req.NodeId,
			NodeName:  result.NodeName,
			Step:      0,
			Error:     result.Error,
			Duration:  result.Duration,
			Status:    spaceRecord.Status,
			StartTime: utils.FormatDate(startTime),
		},
	}
	if err != nil {
		logx.Errorw("[画布] 获取节点执行结果失败",
			logx.Field("工作空间ID", req.Id),
			logx.Field("节点ID", req.NodeId),
			logx.Field("错误", err))
		return nil, errors.New(int(logic.SystemError), "获取节点执行结果失败")
	}
	return resp, nil
}
