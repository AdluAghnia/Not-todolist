package handler

import (
	"log"

	"github.com/AdluAghnia/not_todolist/database"
	"github.com/AdluAghnia/not_todolist/middleware"
	"github.com/AdluAghnia/not_todolist/models"
	"github.com/AdluAghnia/not_todolist/repository"
	"github.com/gofiber/fiber/v2"
)

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
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    timePassed := repository.GetTimeSinceCreated(todos)
    return c.Render("todoList", fiber.Map{
        "Tasks": todos,
        "User": *user,
        "TimePassed": timePassed,
    })
}

func UpdateTodoViewHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return err
    }
    todo, err := repository.GetTodoByID(db, c.Params("id"))
    if err != nil {
        return err
    }

    return c.Render("updateForm", fiber.Map{
        "Task": todo,
    }, "layouts/main")
}

// TODO: Completed this function
func UpdateTodoHandler(c *fiber.Ctx) error {
    db, err := database.Db()

    if err != nil {
        return err
    }

    id := c.Params("id")
    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return err
    }

    todo, err := repository.GetTodoByID(db, id)
    if err != nil {
        return err
    }
    todo.Title = c.FormValue("title")
    todo.Description = c.FormValue("description")
    status := c.FormValue("status")

    log.Println(c.FormValue("title"))
    log.Println(c.FormValue("title"))

    if status == "done" {
        todo.Completed = true
    } else {
        todo.Completed = false
    }

    if err := db.Save(&todo).Error; err != nil {
        return err
    }

    todos, err := repository.GetTodosByID(db, int(user.ID))
    if err != nil {
        return err
    }
    
    return c.Render("todoList", fiber.Map{
        "Tasks": todos,
    }, "layouts/main")
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


