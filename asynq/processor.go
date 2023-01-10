package asynq

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bsidhu009/go-task-queue/async/task"
	"github.com/bsidhu009/go-task-queue/log"
)

type Processor struct {
	logger  *log.Logger
	handler Handler
	// channel to communicate back to the long-running "processor" goroutine.
	// once is used to send value to the channel only once.
	done chan struct{}

	// quit channel is closed when the shutdown of the "processor" goroutine starts.
	quit chan struct{}

	// abort channel communicates to the in-flight worker goroutines to stop.
	abort chan struct{}
}

type ProcessorParams struct {
	logger *log.Logger

	queues map[string]int
}

func NewProcessor(params ProcessorParams) *Processor {
	return &Processor{
		done:    make(chan struct{}),
		quit:    make(chan struct{}),
		abort:   make(chan struct{}),
		handler: HandlerFunc(func(ctx context.Context, t *task.Task) error { return fmt.Errorf("handler not set") }),
	}
}

func (p *Processor) start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-p.done:
				p.logger.Debug("Processor done")
				return
			default:
				p.exec()
			}
		}
	}()
}

type TaskMessage struct {
	// Type indicates the kind of the task to be performed.
	Type string

	// Payload holds data needed to process the task.
	Payload []byte

	// ID is a unique identifier for each task.
	ID string

	// Queue is a name this message should be enqueued to.
	Queue string
}

func (p *Processor) exec() {
	resCh := make(chan error, 1)

	ctx := context.Background()
	msg := TaskMessage{
		Type:    "type",
		Payload: nil,
		ID:      "id",
		Queue:   "queue:id",
	}

	go func() {
		task := task.NewTask(
			msg.Type,
			msg.Payload,
		)
		resCh <- p.perform(ctx, task)
	}()
}

func (p *Processor) perform(ctx context.Context, task *task.Task) (err error) {
	return p.handler.ProcessTask(ctx, task)
}

var SkipRetry = errors.New("skip retry for the task")
