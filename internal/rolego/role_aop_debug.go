package rolego

import (
	"encoding/json"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/engine"
	"github.com/zeromicro/go-zero/core/logx"
)

// 链路追踪 AOP
var (
	// before around-before after around-after
	// 存储日志
	// _ types.BeforeAspect = (*DebugAop)(nil)
	_ types.AfterAspect = (*DebugAop)(nil)
	// 链路追踪
	_ types.AroundAspect = (*DebugAop)(nil)
)

// DebugAop 节点Trace日志切面
type DebugAop struct {
}

// logx.Infof("debug Around before ruleChainId:%s,nodeName:%s,nodeType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Type, ctx.Self().GetNodeId().Id, msg, relationType)

// Around 计算耗时
func (aspect *DebugAop) Around(ctx types.RuleContext, msg types.RuleMsg, relationType string) (types.RuleMsg, bool) {
	inputMar, _ := json.MarshalIndent(msg, "", "    ")
	logx.Infof("around-before:%s,%s, %s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, inputMar)

	// 执行当前节点
	ctx.Self().OnMsg(ctx, msg)
	return msg, false
}

func (aspect *DebugAop) After(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	// output
	outputMar, _ := json.MarshalIndent(msg, "", "    ")
	logx.Infof("around-output: %s,%s,%s", ctx.RuleChain().GetNodeId().Id, ctx.Self().(*engine.RuleNodeCtx).SelfDefinition.Name, outputMar)
	if err != nil {
		logx.Errorf("error: %s", err.Error())
	}
	return msg
}

func (aspect *DebugAop) New() types.Aspect {
	return &DebugAop{}
}

// Order 值越小越优先执行
func (aspect *DebugAop) Order() int {
	return 900
}

// PointCut 切入点 所有节点都会执行
func (aspect *DebugAop) PointCut(ctx types.RuleContext, msg types.RuleMsg, relationType string) bool {
	return true
}
