package users

import (
	"github.com/aftaab60/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	StatusActive = "active"
)

type User struct {
	Id int64 `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	DateCreated string `json:"dateCreated"`
	Password string `json:"password"`
	Status string `json:"status"`
}

type Users []User

func (user *User) Validate() *rest_errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	if user.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	return nil
}
