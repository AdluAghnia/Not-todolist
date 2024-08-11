package routes

import (
	"github.com/AdluAghnia/not_todolist/handler"
	"github.com/AdluAghnia/not_todolist/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    // Auth Routes
    app.Get("/login", handler.ViewLogin)
    app.Get("/register", handler.ViewRegister)
    app.Post("/api/login", handler.LoginHandler)
    app.Post("/api/register", handler.RegisterHandler)
    app.Get("/logout", handler.LogoutHandler) 

    app.Get("/",middleware.JWTMiddleware() , handler.IndexHandler)
    app.Get("/todo", middleware.JWTMiddleware(), handler.ViewAddTask)
    app.Get("/userinfo", middleware.JWTMiddleware(), handler.GetUserInformation)

   app.Post("/api/todo", middleware.JWTMiddleware(), handler.AddTaskHandler)
}
