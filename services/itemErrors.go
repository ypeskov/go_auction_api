package services

import "fmt"

type ItemError struct {
	Code    string
	Message string
}

func NewItemError(code, message string) *ItemError {
	return &ItemError{
		Code:    code,
		Message: message,
	}
}

func (e ItemError) Error() string {
	return fmt.Sprintf("Code: %s. Message: %s", e.Code, e.Message)
}

var IncorrectUserRoleErr = ItemError{
	Code:    "INCORRECT_USER_ROLE",
	Message: "user must be a seller",
}
