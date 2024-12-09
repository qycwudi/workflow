package rulego

import (
	"encoding/json"
	"time"

	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	enums "workflow/internal/types"
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
	if msg.Type == "CANVAS_MSG" {
		// 画布运行
		logx.Infof("CANVAS START ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Start", ctx.Self().GetNodeId().Id, msg)
		_, err := RoleChain.svc.SpaceRecordModel.Insert(ctx.GetContext(), &model.SpaceRecord{
			WorkspaceId:  ctx.RuleChain().GetNodeId().Id,
			Status:       enums.RecordStatusRunning,
			SerialNumber: msg.Id,
			RunTime:      time.Now(),
			RecordName:   msg.Id,
		})
		if err != nil {
			logx.Errorf("create space record err:%s", err.Error())
		}
	} else {
		// API 调用
		logx.Infof("API START ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Start", ctx.Self().GetNodeId().Id, msg)
		msgMar, _ := json.Marshal(msg)
		_, err := RoleChain.svc.ApiRecordModel.Insert(ctx.GetContext(), &model.ApiRecord{
			Status:   enums.RecordStatusRunning,
			TraceId:  msg.Id,
			Param:    string(msgMar),
			Extend:   "{}",
			CallTime: time.Now(),
		})
		if err != nil {
			logx.Errorf("create api record err:%s", err.Error())
		}
	}
	return msg
}

func (aspect *RunAop) End(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	status := enums.RecordStatusSuccess
	if err != nil {
		status = enums.RecordStatusFail
		logx.Info(err.Error())
	}
	if msg.Type == "CANVAS_MSG" {
		logx.Infof("CANVAS END ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "End", ctx.Self().GetNodeId().Id, msg, relationType)
		err = RoleChain.svc.SpaceRecordModel.UpdateStatusBySid(ctx.GetContext(), msg.Id, status)
		if err != nil {
			logx.Errorf("update space record status err:%s", err.Error())
		}
	} else {
		logx.Infof("API END ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "End", ctx.Self().GetNodeId().Id, msg, relationType)
		err = RoleChain.svc.ApiRecordModel.UpdateStatusByTraceId(ctx.GetContext(), msg.Id, status)
		if err != nil {
			logx.Errorf("update api record status err:%s", err.Error())
		}
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
