package auth

import (
	"github.com/AdluAghnia/not_todolist/database"
	"github.com/AdluAghnia/not_todolist/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func ComparePasswordHash(password, hash string) (bool, error) {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        return false, err
    }
    return true, nil
}

func FindUserByEmail(email string) (models.User, error) {
    db, err := database.Db()
    if err != nil {
        return models.User{}, err
    }
    
    user := models.User{}

    db.Where("email = ?", email).Find(&user)

    return user, nil
}
