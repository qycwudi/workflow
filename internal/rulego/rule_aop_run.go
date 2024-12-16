package rulego

import (
	"encoding/json"
	"time"

	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"

	enums "workflow/internal/enum"
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
	if msg.Type == enums.CanvasMsg {
		// 画布运行
		logx.Infof("CANVAS START ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Start", ctx.Self().GetNodeId().Id, msg)
		_, err := RoleChain.svc.SpaceRecordModel.Insert(ctx.GetContext(), &model.SpaceRecord{
			WorkspaceId:  ctx.RuleChain().GetNodeId().Id,
			Status:       enums.TraceStatusRunning,
			SerialNumber: msg.Id,
			RunTime:      time.Now(),
			RecordName:   msg.Id,
			Other:        "{}",
		})
		if err != nil {
			logx.Errorf("create space record err:%s", err.Error())
			ctx.TellFailure(msg, err)
		}
	} else {
		// API 调用
		logx.Infof("API START ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v", ctx.RuleChain().GetNodeId().Id, "Start", ctx.Self().GetNodeId().Id, msg)
		msgMar, _ := json.Marshal(msg)
		_, err := RoleChain.svc.ApiRecordModel.Insert(ctx.GetContext(), &model.ApiRecord{
			Status:     enums.RecordStatusRunning,
			TraceId:    msg.Id,
			Param:      string(msgMar),
			Extend:     "{}",
			CallTime:   time.Now(),
			ApiId:      msg.Metadata["api_id"],
			ApiName:    msg.Metadata["api_name"],
			SecretyKey: msg.Metadata["secret_key"],
			ErrorMsg:   "",
		})
		if err != nil {
			logx.Errorf("create api record err:%s", err.Error())
		}
	}
	return msg
}

func (aspect *RunAop) End(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) types.RuleMsg {
	status := enums.RecordStatusSuccess
	var errMsg string
	if err != nil {
		status = enums.RecordStatusFail
		logx.Info(err.Error())
		errMsg = err.Error()
	}
	if msg.Type == enums.CanvasMsg {
		logx.Infof("CANVAS END ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "End", ctx.Self().GetNodeId().Id, msg, relationType)
		// 获取开始时间
		startTime := msg.Metadata.Values()["startTime"]
		if startTime == "" {
			logx.Errorf("start time is empty")
			startTime = time.Now().Format(time.DateTime)
		}
		start, _ := time.ParseInLocation(time.DateTime, startTime, time.Local)
		duration := time.Since(start).Milliseconds()
		err = RoleChain.svc.SpaceRecordModel.UpdateStatusBySid(ctx.GetContext(), msg.Id, status, duration)
		if err != nil {
			logx.Errorf("update space record status err:%s", err.Error())
			ctx.TellFailure(msg, err)
		}
	} else {
		logx.Infof("API END ruleChainId:%s,flowType:%s,nodeId:%s,msg:%+v,relationType:%s", ctx.RuleChain().GetNodeId().Id, "End", ctx.Self().GetNodeId().Id, msg, relationType)
		err = RoleChain.svc.ApiRecordModel.UpdateStatusByTraceId(ctx.GetContext(), msg.Id, status, errMsg)
		if err != nil {
			logx.Errorf("update api record status err:%s", err.Error())
			ctx.TellFailure(msg, err)
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
