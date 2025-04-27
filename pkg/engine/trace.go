package engine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
)

var Trace *TraceModel

type TraceModel struct {
	TraceModel model.TraceModel
}

func NewTrace(model model.TraceModel) {
	Trace = &TraceModel{
		TraceModel: model,
	}
}

type TraceRecore struct {
	WorkspaceId string    `db:"workspace_id"` // 空间 ID
	TraceId     string    `db:"trace_id"`     // 追踪 ID
	Input       any       `db:"input"`        // 组件输入
	Logic       any       `db:"logic"`        // 执行逻辑
	Output      any       `db:"output"`       // 组件输出
	Step        int64     `db:"step"`         // 分步
	NodeId      string    `db:"node_id"`      // 节点 ID
	NodeName    string    `db:"node_name"`    // 节点名称
	Status      string    `db:"status"`       // 运行状态
	ElapsedTime int64     `db:"elapsed_time"` // 运行耗时
	StartTime   time.Time `db:"start_time"`   // 执行时间
	ErrorMsg    string    `db:"error_msg"`    // 错误信息
}

func (t *TraceModel) CreateTrace(ctx context.Context, trace *TraceRecore) {
	input, err := json.Marshal(trace.Input)
	if err != nil {
		logx.Error(ctx, err)
	}
	var logic string
	if trace.Logic != nil {
		logics, err := json.Marshal(trace.Logic)
		if err != nil {
			logx.Error(ctx, err)
		}
		logic = string(logics)
	} else {
		logic = "{}"
	}
	var output []byte
	if trace.Output != "" {
		output, err = json.Marshal(trace.Output)
		if err != nil {
			logx.Error(ctx, err)
		}
	}
	result, err := t.TraceModel.Insert(ctx, &model.Trace{
		WorkspaceId: trace.WorkspaceId,
		TraceId:     trace.TraceId,
		Input:       string(input),
		Output:      string(output),
		Logic:       logic,
		Step:        trace.Step,
		NodeId:      trace.NodeId,
		NodeName:    trace.NodeName,
		Status:      trace.Status,
		StartTime:   trace.StartTime,
		ElapsedTime: trace.ElapsedTime,
		ErrorMsg:    trace.ErrorMsg,
	})
	if err != nil {
		logx.Error(ctx, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logx.Error(ctx, err)
	}
	logx.Infof("create trace result: %v", rows)
}

func (t *TraceModel) UpdateTrace(ctx context.Context, trace *TraceRecore) {
	var output []byte
	var err error
	if trace.Output != "" {
		output, err = json.Marshal(trace.Output)
		if err != nil {
			logx.Error(ctx, err)
		}
	}
	t.TraceModel.UpdateByTraceIdAndNodeId(ctx, &model.Trace{
		TraceId:     trace.TraceId,
		NodeId:      trace.NodeId,
		Status:      trace.Status,
		Output:      string(output),
		ErrorMsg:    trace.ErrorMsg,
		ElapsedTime: trace.ElapsedTime,
	})
}
