package canvas

import (
	"context"
	json2 "encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type CanvasRunSingleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunSingleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunSingleDetailLogic {
	return &CanvasRunSingleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunSingleDetailLogic) CanvasRunSingleDetail(req *types.CanvasRunSingleDetailRequest) (resp *types.CanvasRunSingleDetailResponse, err error) {
	trace, err := l.svcCtx.TraceModel.FindOneByNodeIdAndWorkspaceId(l.ctx, req.Id, req.NodeId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询运行记录失败")
	}
	// 将字符串转换为map
	var input, output map[string]interface{}
	if err := json2.Unmarshal([]byte(trace.Input), &input); err != nil {
		l.Errorf("解析输入参数失败 err:%+v", err)
		return nil, errors.New(int(logic.SystemError), "解析输入参数失败")
	}
	if err := json2.Unmarshal([]byte(trace.Output), &output); err != nil {
		l.Errorf("解析输出结果失败 err:%+v", err)
		return nil, errors.New(int(logic.SystemError), "解析输出结果失败")
	}

	// 如果data字段是字符串,则尝试解析成map
	if dataStr, ok := input["data"].(string); ok {
		var dataMap map[string]interface{}
		if err := json2.Unmarshal([]byte(dataStr), &dataMap); err == nil {
			input["data"] = dataMap
		}
	}
	if dataStr, ok := output["data"].(string); ok {
		var dataMap map[string]interface{}
		if err := json2.Unmarshal([]byte(dataStr), &dataMap); err == nil {
			output["data"] = dataMap
		}
	}

	// 格式化 input 和 output 为美化的 JSON 字符串
	inputJSON, err := json2.MarshalIndent(input, "", "  ")
	if err != nil {
		l.Errorf("failed to format input params err:%+v", err)
		return nil, errors.New(int(logic.SystemError), "格式化输入参数失败")
	}

	outputJSON, err := json2.MarshalIndent(output, "", "  ")
	if err != nil {
		l.Errorf("failed to format output result err:%+v", err)
		return nil, errors.New(int(logic.SystemError), "格式化输出结果失败")
	}

	resp = &types.CanvasRunSingleDetailResponse{
		NodeId:    trace.NodeId,
		NodeName:  trace.NodeName,
		StartTime: trace.StartTime.UnixMilli(),
		Duration:  trace.ElapsedTime,
		Status:    trace.Status,
		Error:     trace.ErrorMsg,
		Input:     string(inputJSON),
		Output:    string(outputJSON),
	}

	return resp, nil
}
