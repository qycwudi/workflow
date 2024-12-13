package rulego

import (
	"context"
	"encoding/json"
	"time"
	enums "workflow/internal/enum"

	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/engine"

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

// logx.Infof("debug Around before ruleChainId:%s,nodeName:%s,nodeType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Type, ctx.Self().GetNodeId().Id, msg, relationType)

// Around 计算耗时
func (aspect *TraceAop) Around(ctx types.RuleContext, msg types.RuleMsg, relationType string) (types.RuleMsg, bool) {
	start := time.Now() // 记录开始时间
	// input
	inputMar, _ := json.MarshalIndent(msg, "", "    ")
	// logic
	logicMar, _ := json.MarshalIndent(ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Configuration, "", "    ")
	// logx.Infof("around-before:%s,%s, %s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, inputMar)

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
	// logx.Infof("around-耗时: %s,%s,%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, elapsed)

	// output
	outputMar, _ := json.MarshalIndent(ctx.GetContext().Value(ctx.RuleChain().GetNodeId().Id), "", "    ")
	// logx.Infof("around-output: %s,%s,%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, outputMar)
	// 更新追踪
	traceQueue <- &model.Trace{
		TraceId:     msg.Id,
		NodeId:      ctx.Self().GetNodeId().Id,
		ElapsedTime: elapsed.Microseconds(),
		Output:      string(outputMar),
		Status:      enums.TraceStatusFinish,
	}
	return msg, false
}

func (aspect *TraceAop) After(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	ctx.SetContext(context.WithValue(ctx.GetContext(), ctx.RuleChain().GetNodeId().Id, msg))
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
