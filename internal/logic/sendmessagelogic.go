package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gogogo/internal/asynq"
	"gogogo/internal/svc"
	"gogogo/internal/types"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

/*
	docker run -d --name rocketmq-dashboard -e "JAVA_OPTS=-Drocketmq.namesrv.addr=10.8.0.61:9876" -p 8088:8080 -t apacherocketmq/rocketmq-dashboard:latest
	/Users/qiangyuecheng/rockerMQ
	docker-compose up -0
	docker-compose down
*/

func (l *SendMessageLogic) SendMessage(req *types.SendMessageRequest) (resp *types.SendMessageResponse, err error) {
	resp = &types.SendMessageResponse{Code: int(SUCCESS), Message: "SUCCESS"}
	// 分发消息
	if req.NeedOcr {
		// 发送ocr消息
		err := asynq.SendOcrMessage(l.ctx, l.svcCtx.AsynqTaskClient, asynq.Message2OcrPayload(req))
		if err != nil {
			resp.Code = int(SendMessageErr)
			resp.Message = "发送OCR消息异常"
		}
		return resp, err
	}

	if req.NeedLlm {
		// 发送特征提取消息
		err := asynq.SendLlmMessage(l.ctx, l.svcCtx.AsynqTaskClient, asynq.Message2LlmPayload(req))
		if err != nil {
			resp.Code = int(SendMessageErr)
			resp.Message = "发送LLM消息异常"
		}
		return resp, err
	}

	// 目前都需要llm,后续如果接入更多的流程请修改这里的错误提示 ocr->llm
	l.Errorf("no need to deal with :%+v", req)
	resp.Code = int(SendMessageParamFormattingError)
	resp.Message = "发送消息格式错误,needLlm is false"
	return resp, nil
}
