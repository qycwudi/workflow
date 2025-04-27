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
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}

	err = workflow.Register(l.ctx, canvas.Draft)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "注册任务流失败")
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
		Id: traceId,
	}
	go run(context.Background(), l.svcCtx, traceId, canvas.WorkspaceId, data)
	return resp, nil
}

func run(ctx context.Context, svcCtx *svc.ServiceContext, traceId, workspaceId string, data map[string]any) {
	// 运行文件
	_, result, err := workflow.Run(ctx, traceId, workspaceId, data)
	if err != nil {
		logx.Errorw("Failed to run the task flow", logx.Field("error", err))
	}
	logx.Infow("Run the task flow successfully", logx.Field("result", result))
	record, err := svcCtx.SpaceRecordModel.FindOneBySerialNumber(ctx, traceId)
	if err != nil {
		logx.Errorw("Failed to find the running record", logx.Field("error", err))
	}

	other, _ := json.Marshal(result.Output)
	record.Other = string(other)
	status := enum.RecordStatusSuccess
	if result.Error != "" {
		status = enum.RecordStatusFail
		other, _ = json.Marshal(map[string]string{
			"error": result.Error,
		})
	}
	record.Status = status

	record.Duration = time.Since(record.RunTime).Milliseconds()

	err = svcCtx.SpaceRecordModel.Update(ctx, record)
	if err != nil {
		logx.Errorw("更新运行记录失败", logx.Field("error", err))
	}
}
