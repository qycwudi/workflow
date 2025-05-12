package trace

import (
	"context"
	"sort"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/enum"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type QueryComponentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询组件运行结果
func NewQueryComponentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryComponentsLogic {
	return &QueryComponentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryComponentsLogic) QueryComponents(req *types.QueryComponentsRequest) (resp *types.QueryComponentsResponse, err error) {
	// 完成
	resp = &types.QueryComponentsResponse{
		Status:  string(enum.RecordStatusFail),
		Records: make([]types.Record, 0),
	}
	// 查询画布运行记录
	record, err := l.svcCtx.SpaceRecordModel.FindBySerialNumber(l.ctx, req.SerialId)
	if err != nil {
		l.Errorf("query components space record err:%+v", err)
		return resp, errors.New(int(logic.SystemOrmError), "查询画布运行记录失败")
	}

	resp.Status = record.Status

	// 查询组件运行结果
	traces, err := l.svcCtx.TraceModel.FindByTraceId(l.ctx, req.SerialId)
	if err == model.ErrNotFound {
		return resp, errors.New(int(logic.SystemOrmError), "没有找到组件运行结果")
	}
	if err != nil {
		l.Errorf("query components trace err:%+v", err)
		return resp, errors.New(int(logic.SystemOrmError), "查询组件运行结果失败")
	}

	// 转换组件运行结果
	records := make([]types.Record, len(traces))
	for i, trace := range traces {
		records[i] = types.Record{
			Input:     trace.Input,
			Output:    trace.Output,
			NodeId:    trace.NodeId,
			NodeType:  trace.NodeType,
			NodeName:  trace.NodeName,
			Step:      trace.Step,
			Error:     trace.ErrorMsg,
			Duration:  trace.ElapsedTime,
			Status:    trace.Status,
			StartTime: utils.FormatDate(trace.StartTime),
		}
	}
	// 根据 step 升序
	sort.Slice(records, func(i, j int) bool {
		return records[i].Step < records[j].Step
	})
	resp.Records = records
	return resp, nil
}
