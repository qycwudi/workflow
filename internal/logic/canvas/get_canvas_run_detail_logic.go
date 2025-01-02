package canvas

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetCanvasRunDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCanvasRunDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCanvasRunDetailLogic {
	return &GetCanvasRunDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCanvasRunDetailLogic) GetCanvasRunDetail(req *types.GetCanvasRunDetailReq) (resp *types.GetCanvasRunDetailResp, err error) {

	traces, err := l.svcCtx.TraceModel.FindByTraceId(l.ctx, req.RecordId)
	if err != nil {
		l.Errorf("GetCanvasRunDetail err:%+v", err)
		return nil, errors.New(int(logic.SystemOrmError), "查询运行记录详情失败")
	}

	// 查询运行记录
	record, err := l.svcCtx.SpaceRecordModel.FindBySerialNumber(l.ctx, traces[0].TraceId)
	if err != nil {
		l.Errorf("GetCanvasRunRecord err:%+v", err)
		return nil, errors.New(int(logic.SystemOrmError), "查询运行记录失败")
	}

	// 转换记录格式
	components := make([]types.ComponentDetail, 0, len(traces))
	for _, trace := range traces {
		// 将字符串转换为map
		var input, output map[string]interface{}
		if err := json.Unmarshal([]byte(trace.Input), &input); err != nil {
			l.Errorf("parse input json err:%v", err)
			input = make(map[string]interface{})
		}
		if err := json.Unmarshal([]byte(trace.Output), &output); err != nil {
			l.Errorf("parse output json err:%v", err)
			output = make(map[string]interface{})
		}

		components = append(components, types.ComponentDetail{
			Id:        trace.NodeId,
			Name:      trace.NodeName,
			Logic:     trace.Logic,
			StartTime: trace.StartTime.UnixMilli(),
			Duration:  trace.ElapsedTime,
			Status:    trace.Status,
			Input:     input,
			Output:    output,
			Error:     trace.ErrorMsg,
		})
	}

	resp = &types.GetCanvasRunDetailResp{
		Id:         req.RecordId,
		StartTime:  record.RunTime.Format("2006-01-02 15:04:05"),
		Duration:   record.Duration,
		Status:     record.Status,
		Components: components,
	}

	return
}
