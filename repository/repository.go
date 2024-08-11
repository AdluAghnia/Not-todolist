package repository

import (
	"errors"

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
