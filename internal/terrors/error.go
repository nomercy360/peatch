package terrors

import (
	"net/http"
)

type Error struct {
	Code    int
	E       error
	Message string
}

func (e *Error) Error() string {
	return e.E.Error()
}

func NotFound(err error) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		E:       err,
		Message: "not found",
	}
}

func BadRequest(err error) *Error {
	return &Error{
		Code:    http.StatusBadRequest,
		E:       err,
		Message: "bad request",
	}
}

func InternalServerError(err error) *Error {
	return &Error{
		Code:    http.StatusInternalServerError,
		E:       err,
		Message: "internal server error",
	}
}

func Unauthorized(err error, message string) *Error {
	return &Error{
		Code:    http.StatusUnauthorized,
		E:       err,
		Message: message,
	}
}
