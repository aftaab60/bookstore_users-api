package mysql_utils

import (
	"errors"
	"github.com/aftaab60/bookstore_utils-go/rest_errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr{
	sqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return rest_errors.NewNotFoundError("no record matching given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", errors.New("database error"))
	}

	switch sqlError.Number {
	case 1062:
		return rest_errors.NewBadRequestError("invalid data")
	}
	return rest_errors.NewInternalServerError("error processing request", errors.New("database error"))
}
