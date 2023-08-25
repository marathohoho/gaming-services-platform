package models

import (
	"errors"
	"net/mail"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (u *User) ValidateUserRequest() error {
	if len(u.Email) <= 0 {
		return errors.New("user email is required")
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return errors.New("invalid email address")
	}

	if u.Age < 21 {
		return errors.New("user must be at least 21 years old to register")
	}

	if len(u.FirstName) <= 0 {
		return errors.New("user firstName is required")
	}

	return nil
}
