package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	model "gogogo/internal/model/mongo"
	"gogogo/internal/types"
	"log"
	"strings"
	"time"
)

type OcrPayload struct {
	Key        string   `json:"key"`
	Value      string   `json:"value"`
	SpiderName string   `json:"spiderName"`
	NeedOcr    bool     `json:"needOcr"`
	OcrAdds    []string `json:"ocrAdds"`
	NeedLlm    bool     `json:"needLlm"`
	LlmType    string   `json:"llmType"`
	OcrResult  string   `json:"ocrResult"`
}

func HandleOcrTask(ctx context.Context, t *asynq.Task) error {
	var p OcrPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logx.WithContext(ctx).Infof("execute ocr task: key = %s,value = %v", p.Key, p.OcrAdds)
	// 识别内容
	logx.WithContext(ctx).Info("Identify content")
	var ocrResult = strings.Builder{}
	for _, addr := range p.OcrAdds {
		ocrResult.WriteString(p.ocrApiV1(addr))
	}
	logx.WithContext(ctx).Info(ocrResult.String())
	// 存储结果
	logx.WithContext(ctx).Info("Store ocr results")
	// 更新文档
	data := model.Data{
		Key:       p.Key,
		OcrResult: ocrResult.String(),
	}
	_, err := AsynqTaskContext.MGDataModel.UpdateOcrResultByKey(ctx, &data)
	if err != nil {
		return err
	}
	// next任务
	p.next(ctx)
	return nil
}

func (p OcrPayload) ocrApiV1(url string) string {

	return "处理结果处理结果处理结果处理结果处理结果"
}

func (p OcrPayload) next(ctx context.Context) {
	if p.NeedLlm {
		payload := LlmPayload{
			Key:        p.Key,
			SpiderName: p.SpiderName,
			Value:      p.Value,
			NeedOcr:    p.NeedOcr,
			OcrAdds:    p.OcrAdds,
			NeedLlm:    p.NeedLlm,
			LlmType:    p.LlmType,
		}
		err := SendLlmMessage(ctx, AsynqTaskContext.AsynqTaskClient, payload)
		if err != nil {
			logx.WithContext(ctx).Errorf("ocr task send next task -> llm task error:", err.Error())
		}
	}
}

func Message2OcrPayload(param *types.SendMessageRequest) OcrPayload {
	return OcrPayload{
		Key:        param.Key,
		Value:      param.Value,
		SpiderName: param.SpiderName,
		NeedOcr:    param.NeedOcr,
		OcrAdds:    param.OcrAdds,
		NeedLlm:    param.NeedLlm,
		LlmType:    param.LlmType,
	}
}

func SendOcrMessage(ctx context.Context, client *asynq.Client, ocrPayload OcrPayload) error {
	task, err := NewOCRTask(ocrPayload)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
		return err
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
		return err
	}
	logx.WithContext(ctx).Info("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

func NewOCRTask(ocrProcess OcrPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(ocrProcess)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeOCRRecognize, payload, asynq.Retention(24*time.Hour)), nil
}
