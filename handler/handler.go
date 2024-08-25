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
    return c.Render("index", nil, "layouts/main")
}

func ViewRegister(c *fiber.Ctx) error {
    return c.Render("register-form", fiber.Map{
        "Title": "Register",
    })
}

func ViewLogin(c *fiber.Ctx) error {
    return c.Render("login-form", fiber.Map{
        "Title": "Login",
    })
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
    
    return c.Redirect("/todo", fiber.StatusSeeOther)
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

func GetUserInformation(c *fiber.Ctx) error {
    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(user)
}
