package database

import (
	"github.com/AdluAghnia/not_todolist/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Conncetion to SQLite database
func Db() (*gorm.DB, error){
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    // Auto Migration
    db.AutoMigrate(&models.User{}, &models.Todo{})
    return db, nil
}
