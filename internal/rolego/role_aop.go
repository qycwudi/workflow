package rolego

import (
	"context"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// 保护性检查，接口变动，在编译时中断
var (
	// 存储日志
	_ types.BeforeAspect = (*Trace)(nil)
	_ types.AfterAspect  = (*Trace)(nil)
	// 链路追踪
	_ types.AroundAspect = (*Trace)(nil)

	_ types.StartAspect     = (*Trace)(nil)
	_ types.EndAspect       = (*Trace)(nil)
	_ types.CompletedAspect = (*Trace)(nil)
)

// Trace 节点Trace日志切面
type Trace struct {
}

func (aspect *Trace) Completed(ctx types.RuleContext, msg types.RuleMsg) types.RuleMsg {
	logx.Infof("AOP Completed ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Completed", ctx.Self().GetNodeId().Id, msg)
	return msg
}

func (aspect *Trace) Start(ctx types.RuleContext, msg types.RuleMsg) types.RuleMsg {
	logx.Infof("AOP Start ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Start", ctx.Self().GetNodeId().Id, msg)
	return msg
}

func (aspect *Trace) End(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	if err != nil {
		logx.Info(err.Error())
	}
	logx.Infof("AOP End ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "End", ctx.Self().GetNodeId().Id, msg, relationType)
	return msg
}

func (aspect *Trace) Around(ctx types.RuleContext, msg types.RuleMsg, relationType string) (types.RuleMsg, bool) {
	logx.Infof("debug Around before ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "Around", ctx.Self().GetNodeId().Id, msg, relationType)
	// 执行当前节点
	ctx.Self().OnMsg(ctx, msg)
	logx.Infof("debug Around after ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "Around", ctx.Self().GetNodeId().Id, msg, relationType)
	return msg, false
}

func (aspect *Trace) New() types.Aspect {
	return &Trace{}
}

// Order 值越小越优先执行
func (aspect *Trace) Order() int {
	return 900
}

// PointCut 切入点 所有节点都会执行
func (aspect *Trace) PointCut(ctx types.RuleContext, msg types.RuleMsg, relationType string) bool {
	return true
}

func (aspect *Trace) Before(ctx types.RuleContext, msg types.RuleMsg, relationType string) types.RuleMsg {
	// 记录开始时间
	ctx.SetContext(context.WithValue(ctx.GetContext(), ctx.GetSelfId(), time.Now()))
	aspect.onLog(ctx, types.In, msg, relationType, nil)
	return msg
}

func (aspect *Trace) After(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	startTime, ok := ctx.GetContext().Value(ctx.GetSelfId()).(time.Time)
	if !ok {
		// 如果上下文中没有开始时间，则使用零值
		startTime = time.Time{}
	}

	// 计算耗时
	duration := time.Since(startTime).Milliseconds()
	logx.Infof("Execution time: %d ms", duration)
	aspect.onLog(ctx, types.Out, msg, relationType, nil)
	return msg
}

func (aspect *Trace) onLog(ctx types.RuleContext, flowType string, msg types.RuleMsg, relationType string, err error) {
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
