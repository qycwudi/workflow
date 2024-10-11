package rolego

import (
	"context"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/engine"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
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
	inputMar, _ := json.Marshal(msg)
	logx.Infof("around-before:%s,%s, %s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, inputMar)
	// 执行当前节点
	ctx.Self().OnMsg(ctx, msg)
	elapsed := time.Since(start) // 计算耗时
	logx.Infof("around-耗时: %s,%s,%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, elapsed)
	// input
	outputMar, _ := json.Marshal(ctx.GetContext().Value(ctx.RuleChain().GetNodeId().Id))
	logx.Infof("around-output: %s,%s,%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, outputMar)
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

// insert, err := aspect.Svc.TraceModel.Insert(ctx.GetContext(), &model.Trace{
// Id:          0,
// WorkspaceId: "",
// TraceId:     "",
// Input:       "",
// Logic:       "",
// Output:      "",
// Step:        0,
// NodeId:      "",
// NodeName:    "",
// Status:      "",
// ElapsedTime: 0,
// StartTime:   time.Time{},
// })
