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
	"time"
)

type LlmPayload struct {
	Key        string   `json:"key"`
	Value      string   `json:"value"`
	SpiderName string   `json:"spiderName"`
	NeedOcr    bool     `json:"needOcr"`
	OcrAdds    []string `json:"OcrAdds"`
	NeedLlm    bool     `json:"needLlm"`
	LlmType    string   `json:"llmType"`
}

func HandleLlmTask(ctx context.Context, t *asynq.Task) error {
	var p LlmPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logx.Infof("execute Llm task: key = %s,value = %s", p.Key, p.Value)
	// 识别内容
	logx.WithContext(ctx).Info("Extract feature")
	// 查询ocr内容和源内容
	source, err := AsynqTaskContext.MGDataModel.FindOneByKey(ctx, p.Key)
	if err != nil {
		return err
	}
	if source.OcrResult != "" {
		source.OcrResult = " content by ocr:" + source.OcrResult
	}
	llmResult := llm(source.Source + source.OcrResult)
	// 存储结果
	logx.WithContext(ctx).Info("Store llm results")
	// 更新文档
	data := model.Data{
		Key:       p.Key,
		LLMResult: llmResult,
	}
	_, err = AsynqTaskContext.MGDataModel.UpdateLlmResultByKey(ctx, &data)
	if err != nil {
		return err
	}
	return nil
}

func llm(source string) map[string]string {
	// 调用llm接口
	return map[string]string{"k": "v"}
}

func Message2LlmPayload(param *types.SendMessageRequest) LlmPayload {
	return LlmPayload{
		Key:        param.Key,
		Value:      param.Value,
		SpiderName: param.SpiderName,
		NeedOcr:    param.NeedOcr,
		OcrAdds:    param.OcrAdds,
		NeedLlm:    param.NeedLlm,
		LlmType:    param.LlmType,
	}
}

func SendLlmMessage(ctx context.Context, client *asynq.Client, LlmPayload LlmPayload) error {
	task, err := NewLlmTask(LlmPayload)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
		return err
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
		return err
	}
	logx.WithContext(ctx).Infof("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
	return nil
}

func NewLlmTask(LlmProcess LlmPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(LlmProcess)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeLLMFeatureExtraction, payload, asynq.Retention(24*time.Hour)), nil
}
