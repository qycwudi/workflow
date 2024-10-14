package rolego

import (
	"context"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/utils"
)

func InitRoleChain(svc *svc.ServiceContext) {
	RoleChain = &roleChain{svc: svc}
}

var RoleChain *roleChain

type roleChain struct {
	svc *svc.ServiceContext
}

func (r *roleChain) LoadChain(id string, json []byte) {
	chain, b := rulego.Get(id)
	if b {
		// 重新加载
		err := chain.ReloadSelf(json)
		if err != nil {
			logx.Errorf("reload self role chain %s fail,err:%v\n", id, err)
		}
		logx.Infof("reload self role chain %s success,json: %s \n", id, string(json))
		return
	}
	config := rulego.NewConfig()
	config.Logger = &utils.RoleCustomLog{}
	_, err := rulego.New(
		id,
		json,
		rulego.WithConfig(config),
		types.WithAspects(&TraceAop{}, &RunAop{}))
	if err != nil {
		logx.Errorf("load role chain fail,err:%v\n", err)
		return
	}
	logx.Infof("load %s role chain success,json:%s \n", id, json)
}

func (r *roleChain) getChain(id string) types.RuleEngine {
	chain, b := rulego.Get(id)
	if !b {
		logx.Errorf("get role chain fail,id:%s\n", id)
	}
	return chain
}

func (r *roleChain) Run(id string, metadata map[string]string, data string) types.RuleMsg {
	logx.Infof("id:%s\n metadata:%+v\n data:%s\n", id, metadata, data)
	chain, b := rulego.Get(id)
	if !b {
		logx.Errorf("get role chain fail,id:%s\n", id)
	}
	metaData := types.NewMetadata()
	for k, v := range metadata {
		metaData.PutValue(k, v)
	}
	msg := types.NewMsg(0, "CANVAS_MSG", types.JSON, metaData, data)
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
	}))
	return result
}

var traceQueue = make(chan *model.Trace, 100) // 带缓冲的通道

func asyncTraceWriter() {
	for entry := range traceQueue {
		writeLogEntry(entry)
	}
}

func writeLogEntry(trace *model.Trace) {
	if trace.Status == "RUNNING" {
		// 新增
		_, err := RoleChain.svc.TraceModel.Insert(context.Background(), trace)
		if err != nil {
			logx.Errorf("roleChain create trace info error: %s", err.Error())
		}
	} else {
		// 更新
		err := RoleChain.svc.TraceModel.UpdateByTraceIdAndNodeId(context.Background(), trace)
		if err != nil {
			logx.Errorf("roleChain update trace info error: %s", err.Error())
		}
	}
}

func init() {
	logx.Infov("start trace log async store")
	go asyncTraceWriter()
}
