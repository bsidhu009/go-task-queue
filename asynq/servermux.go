package asynq

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type ServeMux struct {
	mu  sync.RWMutex
	m   map[string]muxEntry
	es  []muxEntry // slice of entries sorted from longest to shortest.
	mws []MiddlewareFunc
}

// ProcessTask dispatches the task to the handler whose
// pattern most closely matches the task type.
func (mux *ServeMux) ProcessTask(ctx context.Context, task *Task) error {
	h, _ := mux.Handler(task)
	return h.ProcessTask(ctx, task)
}

// Handler If there is no registered handler that applies to the task,
// handler returns a 'not found' handler which returns an error.
func (mux *ServeMux) Handler(t *Task) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	h, pattern = mux.match(t.Type())
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	for i := len(mux.mws) - 1; i >= 0; i-- {
		h = mux.mws[i](h)
	}
	return h, pattern
}

// Find a handler on a handler map given a typename string.
// Most-specific (longest) pattern wins.
func (mux *ServeMux) match(typename string) (h Handler, pattern string) {
	return nil, ""
}

type muxEntry struct {
	h       Handler
	pattern string
}

type MiddlewareFunc func(Handler) Handler

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return new(ServeMux)
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if strings.TrimSpace(pattern) == "" {
		panic("asynq: invalid pattern")
	}
	if handler == nil {
		panic("asynq: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("asynq: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
}

// NotFound returns an error indicating that the handler was not found for the given task.
func NotFound(ctx context.Context, task *Task) error {
	return fmt.Errorf("handler not found for task %q", task.Type())
}

// NotFoundHandler returns a simple task handler that returns a ``not found`` error.
func NotFoundHandler() Handler { return HandlerFunc(NotFound) }
