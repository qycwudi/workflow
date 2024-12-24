package rulego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/engine"
	"github.com/zeromicro/go-zero/core/logx"

	enums "workflow/internal/enum"
	"workflow/internal/model"
)

// 链路追踪 AOP
var (
	// before around-before after around-after
	// 存储日志
	// _ types.BeforeAspect = (*TraceAop)(nil)
	_ types.AfterAspect = (*TraceAop)(nil)
	// 链路追踪
	_ types.AroundAspect = (*TraceAop)(nil)
)

// TraceAop 节点Trace日志切面
type TraceAop struct {
}

// Around 计算耗时
func (aspect *TraceAop) Around(ctx types.RuleContext, msg types.RuleMsg, relationType string) (types.RuleMsg, bool) {
	// 1. 检查请求是否被取消
	if canceled := aspect.checkRequestCanceled(ctx, msg); canceled {
		return msg, false
	}

	// 2. 记录开始时间并准备初始追踪数据
	start := time.Now()
	trace := aspect.prepareInitialTrace(ctx, msg)
	traceQueue <- &trace

	// 3. 执行当前节点
	ctx.Self().OnMsg(ctx, msg)
	elapsed := time.Since(start)

	// 4. 更新追踪结果
	aspect.updateTraceResult(ctx, msg, elapsed)

	return msg, false
}

// 检查请求是否被取消
func (aspect *TraceAop) checkRequestCanceled(ctx types.RuleContext, msg types.RuleMsg) bool {
	select {
	case <-ctx.GetContext().Done():
		trace := model.Trace{
			TraceId:   msg.Id,
			Output:    "{}",
			Step:      0,
			NodeId:    ctx.Self().GetNodeId().Id,
			NodeName:  ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name,
			Status:    enums.TraceStatusFinish,
			ErrorMsg:  "request canceled",
			StartTime: time.Now(),
		}
		traceQueue <- &trace
		logx.Debugf("request canceled,trace:%+v", trace)
		return true
	default:
		return false
	}
}

// 准备初始追踪数据
func (aspect *TraceAop) prepareInitialTrace(ctx types.RuleContext, msg types.RuleMsg) model.Trace {
	// 准备输入数据
	inputMap := make(map[string]interface{})
	var inMsgData interface{}
	_ = json.Unmarshal([]byte(msg.Data), &inMsgData)
	inputMap["msg"] = inMsgData
	inputMap["metadata"] = msg.Metadata
	inputMap["msgType"] = msg.Type
	inputMar, _ := json.MarshalIndent(inputMap, "", "    ")

	// 准备逻辑数据
	logicMar, _ := json.MarshalIndent(ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Configuration, "", "    ")

	return model.Trace{
		WorkspaceId: ctx.RuleChain().GetNodeId().Id,
		TraceId:     msg.Id,
		Input:       string(inputMar),
		Logic:       string(logicMar),
		Output:      "{}",
		Step:        0,
		NodeId:      ctx.Self().GetNodeId().Id,
		NodeName:    ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name,
		Status:      enums.TraceStatusRunning,
		ElapsedTime: 0,
		StartTime:   time.Now(),
	}
}

// 更新追踪结果
func (aspect *TraceAop) updateTraceResult(ctx types.RuleContext, msg types.RuleMsg, elapsed time.Duration) {
	outputMap := make(map[string]interface{})
	resultMsg := ctx.GetContext().Value(nodeIdKey(ctx.RuleChain().GetNodeId().Id)).(types.RuleMsg)

	var outMsgData interface{}
	_ = json.Unmarshal([]byte(resultMsg.Data), &outMsgData)
	outputMap["msg"] = outMsgData
	outputMap["metadata"] = resultMsg.Metadata
	outputMap["msgType"] = resultMsg.Type

	outputMar, _ := json.MarshalIndent(outputMap, "", "    ")

	traceQueue <- &model.Trace{
		TraceId:     msg.Id,
		NodeId:      ctx.Self().GetNodeId().Id,
		ElapsedTime: elapsed.Milliseconds(),
		Output:      string(outputMar),
		Status:      enums.TraceStatusFinish,
		ErrorMsg:    msg.Metadata["error"],
	}
}

type nodeIdKey string

func (aspect *TraceAop) After(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	nodeId := ctx.RuleChain().GetNodeId().Id
	if err != nil {
		msg.Metadata["error"] = err.Error()
	}
	msg.Metadata["relationType"] = relationType

	// 保持原有 context 的取消信号
	parentCtx := ctx.GetContext()
	cctx := context.WithValue(parentCtx, nodeIdKey(nodeId), msg)
	ctx.SetContext(cctx)
	return msg
}

func (aspect *TraceAop) New() types.Aspect {
	return &TraceAop{}
}

// Order 值越小越优先执行
func (aspect *TraceAop) Order() int {
	return 900
}

// PointCut 切入点 所有节点都会执行
func (aspect *TraceAop) PointCut(ctx types.RuleContext, msg types.RuleMsg, relationType string) bool {
	return true
}
