package asynq

import "fmt"

type TaskQueue struct {
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{}
}

func (q *TaskQueue) NewTask(taskType string, payload Task) {
	fmt.Printf("task added queue %+v\n", payload)
}

func (q *TaskQueue) AbortTask(taskId string) {
	fmt.Printf("taskId %+v\n", taskId)
}

func (q *TaskQueue) Close() {

}
