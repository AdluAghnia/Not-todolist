package routes

import (
	"github.com/AdluAghnia/not_todolist/handler"
	"github.com/AdluAghnia/not_todolist/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    app.Get("/",middleware.JWTMiddleware() , handler.IndexHandler)
    app.Get("/migrate", handler.Migrate)
    app.Get("/register", handler.ViewRegister)
    app.Get("/todo", handler.ViewAddTask)
    app.Get("/login", handler.ViewLogin)

    app.Post("/api/login", handler.LoginHandler)
    app.Post("/api/register", handler.RegisterHandler)
    app.Post("/api/todo", handler.AddTaskHandler)
}
