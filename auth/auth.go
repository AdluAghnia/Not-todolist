package auth

import (
	"errors"

	"github.com/AdluAghnia/not_todolist/database"
	"github.com/AdluAghnia/not_todolist/models"
	"github.com/AdluAghnia/not_todolist/repository"
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

func ValidateRegisterRequest(u *models.User) (bool, error) {
    email := u.Email
    username := u.Username
    password := u.Password

    db, err := database.Db()

    if err != nil {
        return false, err
    }

    emailExist, err := repository.UserExistByEmail(db, email)
    if err != nil {
        return false, err
    }

    if !emailExist {
        return false, errors.New("email already exist")
    }
    
    if email == "" {
        return false, errors.New("email require")
    }

    if username == "" {
        return false, errors.New("username cannot empty")
    }

    if len(password) <= 6 {
        return false, errors.New("password should atleast have 6 characters")
    }

    return true, nil
}
