package main

import (
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

	// Instantiate fiber routes.
	app := fiber.New()

	app.Get("/add-new-task", func(c *fiber.Ctx) error {
		task := task.Task{ /* PAYLOAD */ }
		taskid := srv.AddTask(task)
		// FIXME: return taskid or error.
		return nil
	})

	app.Get("/task-abort", func(c *fiber.Ctx) error {
		taskid = FIXME // FIXME: get task id from request
		srv.AbortTask(taskid)
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
