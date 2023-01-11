package main

import (
	"context"
	"log"

	"github.com/bsidhu009/go-task-queue/asynq"
	"github.com/bsidhu009/go-task-queue/asynq/task"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// Instantiate a new server with the specified configuration.
	srv := asynq.NewServer(
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.Handle(task.TypeImageResize, task.ImageResizeTask)
	// other handlers go here.

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

	// Read saved tasks and enqueue them on server
	// FIXME

	// Instantiate fiber routes.
	app := fiber.New()

	app.Get("/add-new-task", func(c *fiber.Ctx) error {
		task := task.Task{ /* PAYLOAD */ }
		ctx = context.Background()
		// FIXME: populate ctx with necessary fields
		taskid := srv.AddTask(task, ctx)
		// FIXME: return taskid or error.
		return nil
	})

	app.Get("/task-abort", func(c *fiber.Ctx) error {
		taskid = FIXME // FIXME: get task id from request
		ctx = context.Background()
		// FIXME: populate ctx with necessary fields
		srv.AbortTask(taskid, ctx)
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
