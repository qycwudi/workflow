package rolego

import (
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"workflow/internal/model"
)

// 画布侧运行监控 AOP
var (
	_ types.StartAspect = (*RunAop)(nil)
	_ types.EndAspect   = (*RunAop)(nil)
)

// RunAop 运行记录
type RunAop struct {
}

func (aspect *RunAop) Start(ctx types.RuleContext, msg types.RuleMsg) types.RuleMsg {
	logx.Infof("AOP START ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Start", ctx.Self().GetNodeId().Id, msg)
	_, err := RoleChain.svc.SpaceRecordModel.Insert(ctx.GetContext(), &model.SpaceRecord{
		WorkspaceId:  ctx.RuleChain().GetNodeId().Id,
		Status:       "running",
		SerialNumber: msg.Id,
		RunTime:      time.Now(),
		RecordName:   msg.Id,
	})
	if err != nil {
		logx.Errorf("create space record err:%s", err.Error())
	}
	return msg
}

func (aspect *RunAop) End(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	logx.Infof("AOP END ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "End", ctx.Self().GetNodeId().Id, msg, relationType)
	status := "success"
	if err != nil {
		status = "fail"
		logx.Info(err.Error())
	}
	record, err := RoleChain.svc.SpaceRecordModel.FindOneBySerialNumber(ctx.GetContext(), msg.Id)
	if err != nil {
		logx.Errorf("update space find record err:%s", err.Error())
		return msg
	}
	record.Status = status
	err = RoleChain.svc.SpaceRecordModel.Update(ctx.GetContext(), record)
	if err != nil {
		logx.Errorf("update space record status err:%s", err.Error())
	}
	return msg
}

func (aspect *RunAop) New() types.Aspect {
	return &RunAop{}
}

// Order 值越小越优先执行
func (aspect *RunAop) Order() int {
	return 100
}

// PointCut 切入点 所有节点都会执行
func (aspect *RunAop) PointCut(ctx types.RuleContext, msg types.RuleMsg, relationType string) bool {
	return true
}
