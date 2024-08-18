package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/AdluAghnia/not_todolist/models"
	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, id uint) (*models.User, error) {
    var user models.User
    if err := db.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
    var user models.User

    err:= db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }

        return nil, err
    }

    return &user, err
}

func UserExistByEmail(db *gorm.DB, email string) (bool, error) {
    var count int64

    if err := db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
        return false, err
    }

    return count < 0, nil 
}

func GetTodosByID(db *gorm.DB, user_id int) ([]models.Todo, error) {
    var todos []models.Todo

    err := db.Where("user_id", user_id).Find(&todos).Error
    if err != nil {
        return nil, err
    }

    return todos, nil
}

func GetTodoByID(db *gorm.DB, id string) (*models.Todo, error) {
    var todo models.Todo

    err:= db.Where("id = ?", id).First(&todo).Error

    if err != nil {
        return nil, err
    }

    return &todo, err
}

func GetTimeSinceCreated(todos []models.Todo) []string {
    var timePassed []string

    for _, todo := range todos {
        if !todo.Completed {
            // Get time since CreatedAt
            duration := time.Since(todo.CreatedAt)

            // Format the duration to be more human-readable
            hours := int(duration.Hours()) % 24
            days := int(duration.Hours()) / 24
            minutes := int(duration.Minutes()) % 60
            seconds := int(duration.Seconds()) % 60
            timeSince := fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds ago", days, hours, minutes, seconds)

            timePassed = append(timePassed, timeSince)
        }

        timePassed = append(timePassed, "Failed")
    }

    return timePassed
}

