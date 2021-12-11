package users

import (
	"fmt"
	"github.com/aftaab60/bookstore_users-api/datasources/mysql/users_db"
	"github.com/aftaab60/bookstore_users-api/logger"
	"github.com/aftaab60/bookstore_users-api/utils/crypto_utils"
	"github.com/aftaab60/bookstore_users-api/utils/errors"
	"github.com/aftaab60/bookstore_users-api/utils/mysql_utils"
)
//only place where we interact with persistence layer
const (
	queryInsertUser = "INSERT into USERS(first_name, last_Name, email, date_created, password, status) values(?, ?, ?, ?, ?, ?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created, password, status from USERS where id=?;"
	queryUpdateUser = "UPDATE users set first_name=?, last_name=?, email=?, password=?, status=? where id=?;"
	queryDeleteUser             = "DELETE from users where id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status from USERS where status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, password, status from USERS where email=? and status=?;"
)

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return errors.NewInternalServerError("database error")
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last inserted user id", err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	return nil
}

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Password, &user.Status); getErr != nil {
		logger.Error(fmt.Sprintf("error when saving user id %d", user.Id), getErr)
		return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Status, user.Id)
	if updateErr != nil {
		return mysql_utils.ParseError(updateErr)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus() ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.Status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results)==0 {
		return nil, errors.NewNotFoundError("no result with matching status: "+user.Status)
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	password := user.Password
	result := stmt.QueryRow(user.Email, user.Status)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Password, &user.Status); getErr != nil {
		logger.Error(fmt.Sprintf("invalid user credenatials"), getErr)
		return errors.NewNotFoundError("incorrect user credentials")
	}
	if err := crypto_utils.CompareBcryptHashWithPassword(user.Password, password); err != nil {
		return errors.NewNotFoundError("incorrect user credentials")
	}
	return nil
}
