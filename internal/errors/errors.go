package errors

import "fmt"

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %s. Message: %s", e.Code, e.Message)
}

var IncorrectReqBodyErr = Error{
	Code:    "INCORRECT_REQUEST_BODY",
	Message: "Failed to parse request body",
}

var ValidationFailedErr = Error{
	Code:    "VALIDATION_FAILED",
	Message: "Failed to validate parameters",
}

var InvalidParamterErr = Error{
	Code:    "INVALID_PARAMETER",
	Message: "Invalid parameter",
}
