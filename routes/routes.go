package routes

import (
	"github.com/AdluAghnia/not_todolist/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    app.Get("/", handler.IndexHandler)
    app.Get("/migrate", handler.Migrate)
    app.Get("/register", handler.ViewRegister)
    app.Get("/todo", handler.ViewAddTask)

    app.Post("/api/register", handler.RegisterHandler)
    app.Post("/api/todo", handler.AddTaskHandler)
}
