package models

import (
	"gorm.io/gorm"
)

type Todo struct {
    gorm.Model
    Title       string
    Description string
    Completed   bool
    UserID      uint // foreign key referencing the User model
    User        User // belongs to relationship
}
