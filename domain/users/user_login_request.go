package users

import (
	"github.com/aftaab60/bookstore_users-api/utils/errors"
	"strings"
)

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func (request *LoginRequest) Validate() *errors.RestErr {
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	if request.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	if request.Password == "" {
		return errors.NewBadRequestError("password is blank")
	}
	return nil
}
