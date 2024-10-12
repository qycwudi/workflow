package canvas

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
	"workflow/internal/utils"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TraceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTraceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TraceLogic {
	return &TraceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TraceLogic) Trace(req *types.TraceRequest) (resp *types.TraceResponse, err error) {
	traceRecord, err := l.svcCtx.TraceModel.FindByTraceId(l.ctx, req.TraceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询运行追踪记录失败")
	}
	traces := make([]types.Trace, len(traceRecord))
	var totoTime int64 = 0
	for i, trace := range traceRecord {
		traces[i] = types.Trace{
			TraceId:     trace.TraceId,
			NodeId:      trace.NodeId,
			NodeName:    trace.NodeName,
			Status:      trace.Status,
			StartTime:   utils.FormatDate(trace.StartTime),
			ElapsedTime: trace.ElapsedTime,
			Input:       trace.Input,
			Output:      trace.Output,
			Logic:       trace.Logic,
		}
		totoTime += trace.ElapsedTime
	}
	resp = &types.TraceResponse{
		Total:            int64(len(traces)),
		TotalElapsedTime: totoTime,
		Traces:           traces,
	}
	return resp, nil
}
