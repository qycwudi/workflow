package rulego

import (
	"context"
	"fmt"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/enum"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/utils"
)

func InitRoleChain(svc *svc.ServiceContext) {
	config := rulego.NewConfig()
	config.Logger = &utils.RoleCustomLog{}
	opts := []types.RuleEngineOption{
		rulego.WithConfig(config),
	}

	RoleChain = &roleChain{svc: svc, opts: opts}
}

var RoleChain *roleChain

type roleChain struct {
	svc  *svc.ServiceContext
	opts []types.RuleEngineOption
}

// 获取当前节点的父节点
func (r *roleChain) GetParentNode(id string, nodeId string) []string {
	chain, b := rulego.Get(id)
	if !b {
		logx.Errorf("get role chain fail,id:%s", id)
	}
	// 获取规则链定义
	def := chain.Definition()
	// 遍历所有连接关系
	var parentNodes []string
	for _, conn := range def.Metadata.Connections {
		// 如果目标节点是当前节点,添加源节点ID到数组
		if conn.ToId == nodeId {
			parentNodes = append(parentNodes, conn.FromId)
		}
	}
	if len(parentNodes) == 0 {
		return []string{}
	}
	return parentNodes
}

// 加载JOB服务链
func (r *roleChain) LoadJobServiceChain(id string, json []byte) error {
	// 默认开启链路追踪
	return r.LoadChain(id, json, true)
}

// 加载API服务链
func (r *roleChain) LoadApiServiceChain(id string, json []byte) error {
	// 读取环境变量配置
	trace := r.svc.Config.RuleServerTrace
	return r.LoadChain(id, json, trace)
}

// 加载画布服务链
func (r *roleChain) LoadCanvasServiceChain(id string, json []byte) error {
	// 默认开启链路追踪
	return r.LoadChain(id, json, true)
}

func (r *roleChain) LoadChain(id string, json []byte, trace bool) error {
	chain, b := rulego.Get(id)
	if b {
		// 重新加载
		err := chain.ReloadSelf(json)
		if err != nil {
			logx.Errorf("reload self role chain %s fail,err:%v\n", id, err)
			return err
		}
		logx.Infof("reload self role chain %s success,json: %s \n", id, string(json))
		return nil
	}
	if trace {
		r.opts = append(r.opts, types.WithAspects(&TraceAop{}, &RunAop{}))
	} else {
		r.opts = append(r.opts, types.WithAspects(&RunAop{}))
	}
	_, err := rulego.New(
		id,
		json,
		r.opts...,
	)
	if err != nil {
		logx.Errorf("load role chain fail,err:%v\n", err)
		return err
	}

	logx.Infof("load %s role chain success,json:%s \n", id, json)
	return nil
}

func (r *roleChain) Run(id string, metadata map[string]string, data string) types.RuleMsg {
	logx.Infof("id:%s metadata:%+v data:%s", id, metadata, data)
	chain, b := rulego.Get(id)
	if !b {
		logx.Errorf("get role chain fail,id:%s", id)
	}
	metaData := types.NewMetadata()
	for k, v := range metadata {
		metaData.PutValue(k, v)
	}
	msg := types.NewMsg(0, enum.CanvasMsg, types.JSON, metaData, data)
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
	}))
	return result
}

var traceQueue = make(chan *model.Trace, 10000)             // 带缓冲的通道
var spaceRecordQueue = make(chan *model.SpaceRecord, 10000) // 带缓冲的通道
var apiRecordQueue = make(chan *model.ApiRecord, 10000)     // 带缓冲的通道

func asyncTraceWriter() {
	for entry := range traceQueue {
		writeTraceLogEntry(entry)
	}
}

func asyncSpaceRecordWriter() {
	for entry := range spaceRecordQueue {
		writeSpaceRecordLogEntry(entry)
	}
}

func asyncApiRecordWriter() {
	for entry := range apiRecordQueue {
		writeApiRecordLogEntry(entry)
	}
}

func writeTraceLogEntry(trace *model.Trace) {
	if trace.Status == enum.TraceStatusRunning {
		// 新增
		_, err := RoleChain.svc.TraceModel.Insert(context.Background(), trace)
		if err != nil {
			logx.Errorf("roleChain create trace info error: %s", err.Error())
		}
		return
	}
	if trace.Status == enum.TraceStatusFinish {
		// 更新
		err := RoleChain.svc.TraceModel.UpdateByTraceIdAndNodeId(context.Background(), trace)
		if err != nil {
			logx.Errorf("roleChain update trace info error: %s", err.Error())
		}
		return
	}
}

func writeSpaceRecordLogEntry(spaceRecord *model.SpaceRecord) {
	_, err := RoleChain.svc.SpaceRecordModel.Insert(context.Background(), spaceRecord)
	if err != nil {
		logx.Errorf("roleChain create space record info error: %s", err.Error())
	}
}

func writeApiRecordLogEntry(apiRecord *model.ApiRecord) {
	_, err := RoleChain.svc.ApiRecordModel.Insert(context.Background(), apiRecord)
	if err != nil {
		logx.Errorf("roleChain create api record info error: %s", err.Error())
	}
}

func init() {
	logx.Infov("start trace log async store")
	go asyncTraceWriter()
	go asyncSpaceRecordWriter()
	go asyncApiRecordWriter()
}

func LoadJobChain() {
	logx.Info("init role server load job")
	jobs, err := RoleChain.svc.JobModel.FindByOn(context.Background())
	if err != nil {
		logx.Errorf("find job server error: %s\n", err.Error())
		return
	}
	successCount := 0
	failCount := 0
	for _, job := range jobs {
		logx.Infof("load job id:%s,name:%s", job.JobId, job.JobName)
		err := RoleChain.LoadJobServiceChain(job.JobId, []byte(job.Dsl))
		if err != nil {
			logx.Errorf("load job %s error: %s\n", job.JobId, err.Error())
			failCount++
		} else {
			successCount++
		}
	}
	logx.Infof("init role server load job complete : total=%d, success=%d, fail=%d", len(jobs), successCount, failCount)
	fmt.Println("init role server load job complete")
}

func LoadApiChain() {
	logx.Info("init role server load api")
	apis, err := RoleChain.svc.ApiModel.FindByOn(context.Background())
	if err != nil {
		logx.Errorf("find api server error: %s\n", err.Error())
		return
	}
	successCount := 0
	failCount := 0
	for _, api := range apis {
		logx.Infof("load api id:%s,name:%s", api.ApiId, api.ApiName)
		err := RoleChain.LoadApiServiceChain(api.ApiId, []byte(api.Dsl))
		if err != nil {
			logx.Errorf("load api %s error: %s\n", api.ApiId, err.Error())
			failCount++
		} else {
			successCount++
		}
	}
	logx.Infof("init role server load api complete : total=%d, success=%d, fail=%d", len(apis), successCount, failCount)
	fmt.Println("init role server load api complete")
}
