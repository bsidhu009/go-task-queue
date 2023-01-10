package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bsidhu009/go-task-queue/asynq"
	"log"
)

type ArchivePayload struct {
	SourceURL string
}

// ArchiveProcessor implements asynq.Handler interface.
type ArchiveProcessor struct {
	// ... fields for struct
}

func NewArchiveProcessor() *ArchiveProcessor {
	return &ArchiveProcessor{}
}

func (processor *ArchiveProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ArchivePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: src=%s", p.SourceURL)
	// Archive resizing code ...
	return nil
}
