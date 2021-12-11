package users

import (
	"encoding/json"
)

type PublicUser struct {
	Id          int64  `json:"id"`
	DateCreated string `json:"dateCreated"`
	Status      string `json:"status"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	DateCreated string `json:"dateCreated"`
	Status      string `json:"status"`
}

func (user *User) Marshall(isPublic bool) interface{} {
	userJson, _ := json.Marshal(user)
	if !isPublic {
		var privateUser PrivateUser
		json.Unmarshal(userJson, &privateUser)
		return privateUser
	}
	var publicUser PublicUser
	json.Unmarshal(userJson, &publicUser)
	return publicUser
}

func (users Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshall(isPublic)
	}
	return result
}
