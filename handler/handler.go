package handler

import (
	"log"
	"time"

	"github.com/AdluAghnia/not_todolist/auth"
	"github.com/AdluAghnia/not_todolist/database"
	"github.com/AdluAghnia/not_todolist/middleware"
	"github.com/AdluAghnia/not_todolist/models"
	"github.com/AdluAghnia/not_todolist/repository"
	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
    return c.SendString("Hello")
}

func ViewRegister(c *fiber.Ctx) error {
    return c.Render("register", fiber.Map{
        "Title": "Register",
    }, "layouts/main")
}

func ViewAddTask(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return err
    }

    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return err
    }

    todos, err := repository.GetTodosByID(db, int(user.ID))
    if err != nil {
        return err
    }

    timePassed := repository.GetTimeSinceCreated(todos)

    return c.Render("todo", fiber.Map{
        "Title": "Task",
        "Tasks": todos,
        "User": *user,
        "TimePassed": timePassed,
    },"layouts/main")
}

func ViewLogin(c *fiber.Ctx) error {
    return c.Render("login", fiber.Map{
        "Title": "Login",
    }, "layouts/main")
}

func LoginHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        log.Printf("Error : %v \n", err.Error())
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    email := c.FormValue("email")
    password := c.FormValue("password")

    if email != "" && password != ""{
        log.Println("OK")
    }

    // Find User by Email
    user, err := repository.GetUserByEmail(db, email)
    if err != nil {
        log.Printf("Error : %v \n", err.Error())
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    // Compare Input Password and stored password
    ok, err := auth.ComparePasswordHash(password, user.Password)
    if !ok && err != nil {
        log.Printf("Error : %v \n", err.Error())
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    // Generate JWT Token
    token, err := middleware.GenerateJWT(user)
    if err != nil {
        log.Printf("Error : %v \n", err.Error())
        return c.Status(fiber.StatusInternalServerError).SendString("Couldn't generate token")
    }

    // Create Cookie
    c.Cookie(&fiber.Cookie{
        Name: "jwt",
        Value: token,
        Expires: time.Now().Add(24 * time.Hour),
    })

    return c.Render("userinfo", fiber.Map{
        "User": user,
    }, "layouts/main")
}

func LogoutHandler(c *fiber.Ctx) error {
    c.Cookie(&fiber.Cookie{
        Name: "jwt",
        Value: "",
        Expires: time.Now().Add(-time.Hour),
        HTTPOnly: true,
        Path: "/",
    })

    return c.Redirect("/", fiber.StatusSeeOther)
}

func AddTaskHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    
    if err != nil {
        return err
    }

    title := c.FormValue("title")
    description := c.FormValue("description")

    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return err
    }

    db.Create(&models.Todo{
        Title: title,
        Description: description,
        Completed: false,
        User: *user,
        UserID: user.ID,
    })

    todos, err := repository.GetTodosByID(db, int(user.ID))
    timePassed := repository.GetTimeSinceCreated(todos)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }
    return c.Render("todoList", fiber.Map{
        "Tasks": todos,
        "User": *user,
        "TimePassed": timePassed,
    })
}

func RegisterHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return c.Render("register", fiber.Map{
            "Error": err.Error(),
        }, "layouts/main")
    }

    email := c.FormValue("email")
    username := c.FormValue("username")
    password := c.FormValue("password")

    user := models.User{
        Username: username,
        Email: email,
        Password: password,
    }

    valid, errMessage := auth.ValidateRegisterRequest(&user)
    
    if !valid {
        return c.Render("register", fiber.Map{
            "Errors": errMessage,
        })
    }
    
    hash, err := auth.HashPassword(user.Password)

    if err != nil {
        return c.Render("register", fiber.Map{
            "Error": err.Error(),
        }, "layouts/main")
    }

    db.Create(&models.User{
        Email: email,
        Username: username,
        Password: hash,
    })

    if c.Get("HX-Request") == "true" {
        c.Set("HX-Redirect", "/login")
        return nil
    }

    return c.Redirect("/login", fiber.StatusSeeOther)
}


func DeleteTodoHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    err = db.Where("id = ?", c.Params("id")).Delete(&models.Todo{}).Error
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
   }

    todos, err := repository.GetTodosByID(db, int(user.ID))
    timePassed := repository.GetTimeSinceCreated(todos)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.Render("todoList", fiber.Map{
        "Tasks": todos,
        "User": *user,
        "TimePassed": timePassed,
    }, "layouts/main")
}

func GetUserInformation(c *fiber.Ctx) error {
    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(user)
}
