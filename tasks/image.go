package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bsidhu009/go-task-queue/asynq"
	"log"
)

type ImageResizePayload struct {
	SourceURL string
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	// ... fields for struct
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Resizing image: src=%s", p.SourceURL)
	// Image resizing code ...
	return nil
}
