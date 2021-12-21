package services

import (
	"errors"
	"github.com/aftaab60/bookstore_users-api/domain/users"
	"github.com/aftaab60/bookstore_users-api/utils/crypto_utils"
	"github.com/aftaab60/bookstore_users-api/utils/date_utils"
	"github.com/aftaab60/bookstore_utils-go/rest_errors"
)

var (
	UserService userServicesInterface = &userService{}
)

type userServicesInterface interface {
	CreateUser(users.User) (*users.User, *rest_errors.RestErr)
	GetUser(int64) (*users.User, *rest_errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *rest_errors.RestErr)
	DeleteUser(int64) *rest_errors.RestErr
	SearchUser(string) (users.Users, *rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, *rest_errors.RestErr)
}

type userService struct{}

func (s *userService) CreateUser(user users.User) (*users.User, *rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.DateCreated = date_utils.GetNowDBString()
	encryptPassword, err := crypto_utils.GetBcryptHash(user.Password)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("password encryption error", errors.New("password error"))
	}
	user.Password = encryptPassword
	if user.Status == "" {
		user.Status = "active"
	}

	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUser(userId int64) (*users.User, *rest_errors.RestErr) {
	result := &users.User{
		Id: userId,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *userService) UpdateUser(isPatch bool, user users.User) (*users.User, *rest_errors.RestErr) {
	currUser, getErr := s.GetUser(user.Id)
	if getErr != nil {
		return nil, getErr
	}
	if isPatch {
		if user.FirstName != "" {
			currUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currUser.LastName = user.LastName
		}
		if user.Email != "" {
			currUser.Email = user.Email
		}
		if user.Password != "" {
			currUser.Password = user.Password
		}
		if user.Status != "" {
			currUser.Status = user.Status
		}
	} else {
		currUser.FirstName = user.FirstName
		currUser.LastName = user.LastName
		currUser.Email = user.Email
		currUser.Password = user.Password
		currUser.Status = user.Status
	}

	if updateErr := currUser.Update(); updateErr != nil {
		return nil, updateErr
	}
	return currUser, nil
}

func (s *userService) DeleteUser(userId int64) *rest_errors.RestErr {
	currUser, getErr := s.GetUser(userId)
	if getErr != nil {
		return getErr
	}

	if err := currUser.Delete(); err != nil {
		return err
	}
	return nil
}

func (s *userService) SearchUser(status string) (users.Users, *rest_errors.RestErr) {
	dao := &users.User{
		Status: status,
	}
	return dao.FindByStatus()
}

func (s *userService) LoginUser(request users.LoginRequest) (*users.User, *rest_errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	dao := &users.User{
		Email:    request.Email,
		Password: request.Password,
		Status: users.StatusActive,
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
