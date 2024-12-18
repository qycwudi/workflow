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
	// 检查请求是否被取消
	select {
	case <-ctx.GetContext().Done():
		// 请求已被取消,记录取消状态并返回
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
		return msg, false
	default:
	}

	start := time.Now() // 记录开始时间
	// input
	inputMar, _ := json.MarshalIndent(msg, "", "    ")
	// logic
	logicMar, _ := json.MarshalIndent(ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Configuration, "", "    ")

	// 新增追踪 todo 加开关
	trace := model.Trace{
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
	traceQueue <- &trace

	// 执行当前节点
	ctx.Self().OnMsg(ctx, msg)
	elapsed := time.Since(start) // 计算耗时 微秒

	// output
	outputMar, _ := json.MarshalIndent(ctx.GetContext().Value(nodeIdKey(ctx.RuleChain().GetNodeId().Id)), "", "    ")
	// 更新追踪
	traceQueue <- &model.Trace{
		TraceId:     msg.Id,
		NodeId:      ctx.Self().GetNodeId().Id,
		ElapsedTime: elapsed.Microseconds(),
		Output:      string(outputMar),
		Status:      enums.TraceStatusFinish,
		ErrorMsg:    msg.Metadata["error"],
	}
	return msg, false
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
