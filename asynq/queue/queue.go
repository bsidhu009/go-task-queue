package asynq

import (
	"fmt"

	"github.com/bsidhu009/go-task-queue/async/task"
)

type TaskQueue struct {
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{}
}

func (q *TaskQueue) NewTask(taskType string, payload task.Task) {
	fmt.Printf("task added queue %+v\n", payload)
}

func (q *TaskQueue) AbortTask(taskId string) {
	fmt.Printf("taskId %+v\n", taskId)
}

func (q *TaskQueue) Close() {

}
