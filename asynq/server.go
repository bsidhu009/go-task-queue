package asynq

import (
	"context"
	"fmt"
	"github.com/bsidhu009/go-task-queue/async/task"
	"github.com/bsidhu009/go-task-queue/log"
	"sync"
)

// Logger supports logging at various log levels.
type Logger interface {
	// Debug logs a message at Debug level.
	Debug(args ...interface{})

	// Info logs a message at Info level.
	Info(args ...interface{})

	// Warn logs a message at Warning level.
	Warn(args ...interface{})

	// Error logs a message at Error level.
	Error(args ...interface{})

	// Fatal logs a message at Fatal level
	// and process will exit with status set to 1.
	Fatal(args ...interface{})
}

// Config specifies the server's background-task processing behavior.
type Config struct {
	// If set to a zero or negative value, NewServer will overwrite the value
	// to the number of CPUs usable by the current process.
	Concurrency int

	// If unset, default logger is used.
	Logger Logger
}

type Handler interface {
	ProcessTask(context.Context, *task.Task) error
}

type HandlerFunc func(context.Context, *task.Task) error

// ProcessTask calls fn(ctx, task)
func (fn HandlerFunc) ProcessTask(ctx context.Context, task *task.Task) error {
	return fn(ctx, task)
}

type Server struct {
	logger *log.Logger

	// wait group to wait for all goroutines to finish.
	wg        sync.WaitGroup
	processor *Processor
}

func NewServer(cfg Config) *Server {
	logger := log.NewLogger(cfg.Logger)
	logger.SetLevel(log.DebugLevel)

	queues := make(map[string]int)

	processor := NewProcessor(ProcessorParams{
		logger: logger,
		queues: queues,
	})

	return &Server{
		processor: processor,
	}
}

func (srv *Server) Run(handler Handler) error {
	if err := srv.Start(handler); err != nil {
		return err
	}
	return nil
}

func (srv *Server) Start(handler Handler) error {
	if handler == nil {
		return fmt.Errorf("asynq: server cannot run with nil handler")
	}

	srv.processor.handler = handler

	if err := srv.start(); err != nil {
		return err
	}

	//srv.logger.Info("Starting processing")
	srv.processor.start(&srv.wg)

	return nil
}

func (srv *Server) start() error {
	return nil
}
