package canvas

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
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
	var index int64 = 0
	workspaceId := traces[0].WorkspaceId
	subTrace := make(map[string][]*model.Trace, 0)
	for _, trace := range traces {
		if trace.WorkspaceId != workspaceId { // 根据workspaceId查询组件列表
			if _, ok := subTrace[trace.WorkspaceId]; !ok {
				subTrace[trace.WorkspaceId] = make([]*model.Trace, 0)
			}
			subTrace[trace.WorkspaceId] = append(subTrace[trace.WorkspaceId], trace)
		}
	}

	// 转换记录格式
	components := make([]types.ComponentDetail, 0, len(traces))
	for _, trace := range traces {
		index++
		if trace.WorkspaceId != workspaceId {
			continue
		}
		// 将字符串转换为map
		var input, output map[string]interface{}
		if err := sonic.Unmarshal([]byte(trace.Input), &input); err != nil {
			l.Errorf("parse input json err:%v", err)
			input = make(map[string]interface{})
		}
		if err := sonic.Unmarshal([]byte(trace.Output), &output); err != nil {
			l.Errorf("parse output json err:%v", err)
			output = make(map[string]interface{})
		}
		subComponents := make([]types.ComponentDetail, 0)

		if subTraces, ok := subTrace[trace.NodeId]; ok {
			for _, subTrace := range subTraces {
				index++
				var subInput, subOutput map[string]interface{}
				if err := sonic.Unmarshal([]byte(subTrace.Input), &subInput); err != nil {
					l.Errorf("parse subTrace input json err:%v", err)
					subInput = make(map[string]interface{})
				}
				if err := sonic.Unmarshal([]byte(subTrace.Output), &subOutput); err != nil {
					l.Errorf("parse subTrace output json err:%v", err)
					subOutput = make(map[string]interface{})
				}
				subComponents = append(subComponents, types.ComponentDetail{
					Id:        subTrace.NodeId,
					Index:     index,
					Name:      subTrace.NodeName,
					Logic:     subTrace.NodeType,
					StartTime: subTrace.StartTime.UnixMilli(),
					Duration:  subTrace.ElapsedTime,
					Status:    subTrace.Status,
					Input:     subInput,
					Output:    subOutput,
					Error:     subTrace.ErrorMsg,
				})
			}
		}

		components = append(components, types.ComponentDetail{
			Id:         trace.NodeId,
			Index:      index,
			Name:       trace.NodeName,
			Logic:      trace.Logic,
			StartTime:  trace.StartTime.UnixMilli(),
			Duration:   trace.ElapsedTime,
			Status:     trace.Status,
			Input:      input,
			Output:     output,
			Error:      trace.ErrorMsg,
			Components: subComponents,
		})
	}

	resp = &types.GetCanvasRunDetailResp{
		Id:         req.RecordId,
		Components: components,
	}

	return
}
