package rolego

import (
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
	"workflow/internal/svc"
)

// 画布侧运行监控 AOP
var (
	_ types.StartAspect = (*RunAop)(nil)
	_ types.EndAspect   = (*RunAop)(nil)
)

// RunAop 运行记录
type RunAop struct {
	svc *svc.ServiceContext
}

func (aspect *RunAop) Start(ctx types.RuleContext, msg types.RuleMsg) types.RuleMsg {
	logx.Infof("AOP Start ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Start", ctx.Self().GetNodeId().Id, msg)

	return msg
}

func (aspect *RunAop) End(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	if err != nil {
		logx.Info(err.Error())
	}
	logx.Infof("AOP End ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "End", ctx.Self().GetNodeId().Id, msg, relationType)
	return msg
}

func (aspect *RunAop) New() types.Aspect {
	return &RunAop{}
}

// Order 值越小越优先执行
func (aspect *RunAop) Order() int {
	return 900
}

// PointCut 切入点 所有节点都会执行
func (aspect *RunAop) PointCut(ctx types.RuleContext, msg types.RuleMsg, relationType string) bool {
	return true
}
