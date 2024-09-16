package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
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

	return nil
}

func NewLlmTask(LlmProcess LlmPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(LlmProcess)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeLLMFeatureExtraction, payload, asynq.Retention(24*time.Hour)), nil
}
