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
    return c.Render("todo", fiber.Map{
        "Title": "Task",
    }, "layouts/main")
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

    var todos []models.Todo
    
    err = db.Find(&todos).Error
    if err != nil {
        return err
    }
    return c.Render("todolist", fiber.Map{
        "Tasks": todos,
        "User": *user,
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

    valid, err := auth.ValidateRegisterRequest(&user)
    
    if !valid && err != nil {
        log.Printf("Error : %v", err.Error())
        return c.Render("list", fiber.Map{
            "Error": err.Error(),
        })
    }
    
    hash, err := auth.HashPassword(user.Password)

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    db.Create(&models.User{
        Email: email,
        Username: username,
        Password: hash,
    })

    var users []models.User
    err = db.Find(&users).Error
    if err != nil {
        return err
    }

    return c.Render("list", fiber.Map{
        "Users": users,
    }, "layouts/main")
}


func GetUserInformation(c *fiber.Ctx) error {
    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(user)
}
