package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email string
    Username string
    Password string
}

func (u *User) ValidateRegister() (bool, error) {
    email := u.Email
    username := u.Username
    password := u.Password

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
