package handler

import (
	"log"

	"github.com/AdluAghnia/not_todolist/database"
	"github.com/AdluAghnia/not_todolist/middleware"
	"github.com/AdluAghnia/not_todolist/models"
	"github.com/AdluAghnia/not_todolist/repository"
	"github.com/gofiber/fiber/v2"
)

func IndexTodoHandler(c *fiber.Ctx) error  {
    db, err := database.Db()
    if err != nil {
        return err
    }

    user, err := middleware.GetUserFromContext(c)
    if err != nil {
        return err
    }

    tasks, err := repository.GetTodosByUserID(db, user.ID)
    if err != nil {
        return err
    }

    if err := repository.UpdateTimeSinceCreated( db, tasks); err != nil {
        log.Println(err.Error())
        return err 
    }

    return c.Render("index-todo", fiber.Map{
        "Tasks": tasks,
    }, "layouts/main")
}

func GetTodoHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return err
    }

    task, err := repository.GetTodoByID(db, c.Params("id"))
    if err != nil {
        return err
    }

    return c.Render("todo", fiber.Map{
        "Task": task,
    }, "layouts/main")
}

func ViewAddTask(c *fiber.Ctx) error {
    return c.Render("todo-form", nil,"layouts/main")
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
    todo := &models.Todo{
        Title: title,
        Description: description,
        Completed: false,
        User: *user,
        UserID: user.ID,
        TimeSinceCreated: "New Task!",
    }

    err = db.Create(todo).Error
    if err != nil {
        return err
    }

    return c.Render("todo", todo, "layouts/main")
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

func UpdateTodoHandler(c *fiber.Ctx) error {
    db, err := database.Db()

    if err != nil {
        return err
    }

    id := c.Params("id")
    todo, err := repository.GetTodoByID(db, id)
    if err != nil {
        return err
    }
    todo.Title = c.FormValue("title")
    todo.Description = c.FormValue("description")
    status := c.FormValue("status")

    if status == "done" {
        todo.Completed = true
        todo.TimeSinceCreated = "failed"
    } else {
        todo.Completed = false
    }


    if err := db.Save(&todo).Error; err != nil {
        return err
    }
    
    todo, err = repository.GetTodoByID(db, id)
    if err != nil {
        return err
    }

    return c.Render("todo", todo, "layouts/main")
}


// TODO : Complete Delete Handler
func DeleteTodoHandler(c *fiber.Ctx) error {
    db, err := database.Db()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    id := c.Params("id")

    err = db.Where("id = ?", id).Delete(&models.Todo{}).Error
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.Send(nil)
}

