package msg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// SyncDraftDslProcessor implements asynq.Handler interface.
/*
	同步草稿DSL
*/
type SyncDraftDslProcessor struct {
	// ... fields for struct
}

type SyncDraftDslPayload struct {
	SourceURL string
}

func (processor *SyncDraftDslProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p SyncDraftDslPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("SyncDraftDsl: src=%s", p.SourceURL)
	// Image resizing code ...
	return nil
}

func NewSyncDraftDslProcessor() *SyncDraftDslProcessor {
	return &SyncDraftDslProcessor{}
}
