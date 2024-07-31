package handler

import (
	"log"
	"github.com/AdluAghnia/not_todolist/database"
	"github.com/AdluAghnia/not_todolist/models"
	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
    return c.SendString("Hello")
}

func Migrate(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return err
    }
    return db.AutoMigrate(&models.User{}, &models.Todo{})
}

func ViewRegister(c *fiber.Ctx) error {
    return c.Render("index", fiber.Map{
        "Title": "Register",
    }, "layouts/main")
}

func ViewAddTask(c *fiber.Ctx) error {
    return c.Render("todo", fiber.Map{
        "Title": "Task",
    }, "layouts/main")
}

func AddTaskHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return err
    }

    title := c.FormValue("title")
    description := c.FormValue("description")

    db.Create(&models.Todo{
        Title: title,
        Description: description,
        Completed: false,
    })

    var todos []models.Todo
    err = db.Find(&todos).Error
    if err != nil {
        return err
    }
    return c.Render("todolist", fiber.Map{
        "Tasks": todos,
    })
}

func RegisterHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return err
    }

    email := c.FormValue("email")
    username := c.FormValue("username")
    password := c.FormValue("password")

    user := models.User{
        Username: username,
        Email: email,
        Password: password,
    }

    valid, err := user.ValidateRegister()
    
    if err != nil {
        log.Printf("Error : %v", err.Error())
        return c.Render("todolist", fiber.Map{
            "Error": err.Error(),
        })
    }

    if valid {
        db.Create(&models.User{
        Email: email,
        Username: username,
        Password: password,
        })
    }

    var users []models.User
    err = db.Find(&users).Error
    if err != nil {
        return err
    }

    return c.Render("list", fiber.Map{
        "Users": users,
    }, "layouts/main")
}
