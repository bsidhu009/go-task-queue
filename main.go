package main

import (
	"github.com/bsidhu009/go-task-queue/asynq"
	"github.com/bsidhu009/go-task-queue/tasks"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	srv := asynq.NewServer(
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.Handle(asynq.TypeImageResize, tasks.NewImageProcessor())

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

	app := fiber.New()

	app.Get("/add-new-task", func(c *fiber.Ctx) error {
		queue := asynq.NewTaskQueue()
		defer queue.Close()

		payload := asynq.Task{}
		queue.NewTask(asynq.TypeArchiveCreate, payload)

		return c.SendString("Hello, World!")
	})

	app.Get("/task-abort", func(c *fiber.Ctx) error {
		queue := asynq.NewTaskQueue()
		defer queue.Close()

		queue.AbortTask("unique-task-1")

		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
