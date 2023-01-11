package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bsidhu009/go-task-queue/asynq"
)

func main() {
	config := asynq.TaskQueueConfig{
		Concurrency: 10,
	}
	taskq, err := asynq.NewTaskQueue(config)
	if err != nil {
		log.Fatal(err)
	}
	go taskq.Run()
	for i := 0; i < 30; i++ {
		var taskType string
		// flip a coin
		if rand.Intn(2) == 0 {
			taskType = asynq.ShortTaskType
		} else {
			taskType = asynq.LongTaskType
		}
		id, err := taskq.AddTaskRequest(taskType, fmt.Sprintf("hello from task number %d", i))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added task %d with type %s\n", id, taskType)
	}
	taskq.Done() // stop when all tasks are done
}
