package rolego

import (
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/engine"
	"github.com/zeromicro/go-zero/core/logx"
)

// 链路追踪 AOP
var (
	// 存储日志
	// _ types.BeforeAspect = (*TraceAop)(nil)
	// _ types.AfterAspect  = (*TraceAop)(nil)
	// 链路追踪
	_ types.AroundAspect = (*TraceAop)(nil)
)

// TraceAop 节点Trace日志切面
type TraceAop struct {
}

func (aspect *TraceAop) Around(ctx types.RuleContext, msg types.RuleMsg, relationType string) (types.RuleMsg, bool) {
	logx.Infof("debug Around before ruleChainId:%s,nodeName:%s,nodeType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Type, ctx.Self().GetNodeId().Id, msg, relationType)
	// 执行当前节点
	ctx.Self().OnMsg(ctx, msg)
	logx.Infof("debug Around after ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "Around", ctx.Self().GetNodeId().Id, msg, relationType)
	return msg, false
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

//
// func (aspect *TraceAop) Before(ctx types.RuleContext, msg types.RuleMsg, relationType string) types.RuleMsg {
// 	// 记录开始时间
// 	ctx.SetContext(context.WithValue(ctx.GetContext(), ctx.GetSelfId(), time.Now()))
// 	aspect.onLog(ctx, types.In, msg, relationType, nil)
// 	return msg
// }
//
// func (aspect *TraceAop) After(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
// 	startTime, ok := ctx.GetContext().Value(ctx.GetSelfId()).(time.Time)
// 	if !ok {
// 		// 如果上下文中没有开始时间，则使用零值
// 		startTime = time.Time{}
// 	}
//
// 	// 计算耗时
// 	duration := time.Since(startTime).Milliseconds()
// 	logx.Infof("Execution time: %d ms", duration)
// 	aspect.onLog(ctx, types.Out, msg, relationType, nil)
// 	return msg
// }

func (aspect *TraceAop) onLog(ctx types.RuleContext, flowType string, msg types.RuleMsg, relationType string, err error) {
	ctx.SubmitTack(func() {
		// 异步记录日志
		if ctx.Self() != nil && ctx.Self().IsDebugMode() {
			var chainId = ""
			if ctx.RuleChain() != nil {
				chainId = ctx.RuleChain().GetNodeId().Id
			}
			// 存储数据库
			logx.Infof("debug ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s,err:%v", chainId, flowType, ctx.Self().GetNodeId().Id, msg, relationType, err)
			// 存储 open observe
		}
	})
}
