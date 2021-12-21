package users

import (
	"github.com/aftaab60/bookstore_utils-go/rest_errors"
	"strings"
)

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func (request *LoginRequest) Validate() *rest_errors.RestErr {
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	if request.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	if request.Password == "" {
		return rest_errors.NewBadRequestError("password is blank")
	}
	return nil
}
