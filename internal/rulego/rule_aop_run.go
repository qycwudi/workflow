package rulego

import (
	"strconv"
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

func (aspect *RunAop) Start(ctx types.RuleContext, msg types.RuleMsg) (types.RuleMsg, error) {
	msg.Metadata[msg.Id+"_param"] = msg.Data
	msg.Metadata["startTime"] = strconv.Itoa(int(time.Now().UnixMilli()))
	logx.Infof("rulego aop run start:%+v", msg)
	return msg, nil
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
		// 获取开始时间
		startTime := msg.Metadata.Values()["startTime"]
		if startTime == "" {
			logx.Errorf("start time is empty")
			startTime = time.Now().Format(time.DateTime)
		}
		start, _ := strconv.Atoi(startTime)
		duration := time.Since(time.UnixMilli(int64(start))).Milliseconds()
		spaceRecordQueue <- &model.SpaceRecord{
			SerialNumber: msg.Id,
			Status:       status,
			Duration:     duration,
			WorkspaceId:  ctx.RuleChain().GetNodeId().Id,
			RunTime:      time.Now(),
			RecordName:   msg.Id,
			Other:        "{}",
		}
	} else {
		apiRecordQueue <- &model.ApiRecord{
			TraceId:    msg.Id,
			Status:     status,
			Extend:     msg.Data,
			Param:      msg.Metadata[msg.Id+"_param"],
			CallTime:   time.Now(),
			ApiId:      msg.Metadata["api_id"],
			ApiName:    msg.Metadata["api_name"],
			SecretyKey: msg.Metadata["secret_key"],
			ErrorMsg:   errMsg,
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
