package openapi

import (
	"context"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/workflow"
)

var (
	recordChan    = make(chan *model.ApiRecord, 10000) // 通道缓冲区大小为10000
	batchSize     = 1000                               // 批量插入的大小
	flushInterval = time.Second                        // 刷新间隔
	once          sync.Once
	stopChan      = make(chan struct{})
	serviceCtx    *svc.ServiceContext
)

// 初始化批量处理
func initBatchProcessor() {
	once.Do(func() {
		go batchProcessor()
	})
}

// 批量处理器
func batchProcessor() {
	var records []*model.ApiRecord
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case record := <-recordChan:
			records = append(records, record)
			if len(records) >= batchSize {
				flushRecords(records)
				records = records[:0]
			}
		case <-ticker.C:
			if len(records) > 0 {
				flushRecords(records)
				records = records[:0]
			}
		case <-stopChan:
			// 处理剩余的记录
			if len(records) > 0 {
				flushRecords(records)
			}
			return
		}
	}
}

// 刷新记录到数据库
func flushRecords(records []*model.ApiRecord) {
	if len(records) == 0 {
		return
	}

	ctx := context.Background()
	err := serviceCtx.ApiRecordModel.BatchInsert(ctx, records)
	if err != nil {
		logx.Errorf("批量插入API记录失败: %v", err)
	}
}

type OpenApiCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type OpenApiCallRequest struct {
	ApiId string         `path:"apiId"` // API ID
	Param map[string]any `json:"param"` // 参数
}

func NewOpenApiCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenApiCallLogic {
	// 确保批量处理器已启动
	initBatchProcessor()
	// 设置全局 ServiceContext
	serviceCtx = svcCtx
	return &OpenApiCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenApiCallLogic) OpenApiCall(req *OpenApiCallRequest) (resp map[string]any, err error) {
	traceId := trace.TraceIDFromContext(l.ctx)
	apiId := req.ApiId
	param := req.Param
	status := "success"
	errMsg := ""
	// 执行
	_, result, err := workflow.Run(l.ctx, traceId, apiId, param)
	if err != nil {
		errMsg = err.Error()
		logx.Errorw("[画布] 执行工作流失败",
			logx.Field("API ID", apiId),
			logx.Field("序列ID", traceId),
			logx.Field("错误", err))
		status = "failed"
	}
	// 记录 api record
	paramBytes, _ := sonic.Marshal(param)
	resultBytes, _ := sonic.Marshal(result)
	record := &model.ApiRecord{
		Status:     status,
		TraceId:    traceId,
		Param:      string(paramBytes),
		Extend:     string(resultBytes),
		CallTime:   time.Now(),
		ApiId:      apiId,
		ApiName:    apiId,
		ErrorMsg:   errMsg,
		SecretyKey: "",
	}
	logx.Infof("API记录: %+v,err:%v", record, err)

	// 将记录发送到通道
	// select {
	// case recordChan <- record:
	// 	// 成功发送到通道
	// default:
	// 	// 通道已满，记录错误
	// 	logx.Errorf("API记录通道已满，丢弃记录: %s", traceId)
	// }

	return result, nil
}
