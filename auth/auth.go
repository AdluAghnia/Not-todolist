package auth

import (
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

func ValidateRegisterRequest(u *models.User) (bool, map[string]string) {
	errorMessage := make(map[string]string)

	email := u.Email
	username := u.Username
	password := u.Password

	db, err := database.Db()

	if err != nil {
		errorMessage["server"] = err.Error()
		return false, errorMessage
	}

	emailExist, err := repository.UserExistByEmail(db, email)

	emailValid := email != ""
	usernameValid := username != ""
	passwordValid := len(password) <= 6

	if err != nil {
		errorMessage["server"] = err.Error()
		return false, errorMessage
	}

	if !emailExist {
		errorMessage["email"] = "Email Already Taken"
	}

	if !emailValid {
		errorMessage["email"] = "Email is required"
	}

	if !usernameValid {
		errorMessage["username"] = "Username required"
	}

	if passwordValid {
		errorMessage["password"] = "password should atleast have 6 characters"
	}

	return emailValid && usernameValid && !passwordValid, errorMessage
}
