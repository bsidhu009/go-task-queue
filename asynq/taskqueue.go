package asynq 

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type TaskQueueConfig struct {
	Concurrency int
}

const MAX_CONCURRENCY = 1000

type TaskQueue struct {
	config   TaskQueueConfig
	taskCh   chan *Task
	taskDone chan int64
	queue    []*Task
	queued   map[int64]*Task
	running  map[int64]*Task
	mu       sync.Mutex
	nextId   int64
	done     bool
}

func NewTaskQueue(config TaskQueueConfig) (*TaskQueue, error) {
	if config.Concurrency < 1 || config.Concurrency > MAX_CONCURRENCY {
		return nil, fmt.Errorf("invalid concurrency: %d", config.Concurrency)
	}
	return &TaskQueue{
		config:   config,
		taskCh:   make(chan *Task, config.Concurrency),
		taskDone: make(chan int64),
		queued:   make(map[int64]*Task),
		running:  make(map[int64]*Task),
		queue:    make([]*Task, 0),
	}, nil
}

func (t *TaskQueue) AddTaskRequest(taskType string, payload interface{}) (int64, error) {
	task, err := CreateTask(taskType, payload)
	if err != nil {
		return -1, err
	}
	(*task).SetId(t.nextId)
	t.nextId++
	t.AddTask(task)
	return (*task).Id(), nil
}

func (t *TaskQueue) AddTask(task *Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	id := (*task).Id()
	if len(t.running) < t.config.Concurrency {
		t.running[id] = task
		t.taskCh <- task
	} else {
		t.queued[id] = task
		t.queue = append(t.queue, task)
	}
}

func (t *TaskQueue) AbortTask(id int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.running[id]; ok {
		// send signal to task to abort
		(*t.running[id]).Abort()
	} else if _, ok := t.queued[id]; ok {
		// remove from queue by iterating over queue.
		// This is inefficient, but ok for a small queue.
		// TBD: To make it more efficient, we could use a skip-list or red-black tree.

		// get position in queue
		pos := -1
		for i, task := range t.queue {
			if (*task).Id() == id {
				pos = i
				break
			}
		}
		if pos == -1 {
			return fmt.Errorf("task %d not found in queue", id)
		}
		t.queue = append(t.queue[:pos], t.queue[pos+1:]...)

		delete(t.queued, id)
	} else {
		return fmt.Errorf("task %d not found", id)
	}
	return nil
}

func (t *TaskQueue) Run() {
	wg := sync.WaitGroup{}
	for i := 0; i < t.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ii := i // capture i
			worker(ii, t.taskCh, t.taskDone)
		}()
	}

	for !t.done {
		done_id := <-t.taskDone
		log.Printf("task %d (%s) done", done_id, (*t.running[done_id]).toString())
		t.mu.Lock()
		delete(t.running, done_id)
		// dispatch next task
		if len(t.queue) > 0 {
			task := t.queue[0]
			t.queue = t.queue[1:]
			id := (*task).Id()
			delete(t.queued, id)
			t.running[id] = task
			t.taskCh <- task
		}
		t.mu.Unlock()
	}
	wg.Wait()
}

func (t *TaskQueue) Done() {
	for {
		// check if all tasks are done, if so, exit. Otherwise, sleep for 1 second
		n := t.Running()
		log.Printf("waiting for %d tasks to finish\n", n)
		if t.Running() == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	t.done = true
	close(t.taskCh)
}

func (t *TaskQueue) Running() int {
	return len(t.running)
}

func worker(workerId int, taskch chan *Task, done chan int64) {
	for {
		task, more := <-taskch
		if !more {
			log.Printf("worker %d exiting\n", workerId)
			return
		}
		(*task).Run()
		done <- (*task).Id()

	}
}
