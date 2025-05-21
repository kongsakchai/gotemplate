package app

import (
	"errors"
	"net/http"
)

type Error struct {
	StatusCd int
	Code     string
	Message  string
	Err      error
}

func NewError(statusCd int, code string, msg string, err ...error) Error {
	var e error
	if len(err) > 0 {
		var errType *Error
		if errors.As(err[0], &errType) {
			return *errType
		}

		e = err[0]
	}

	return Error{
		StatusCd: statusCd,
		Code:     code,
		Message:  msg,
		Err:      e,
	}
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func BadRequestError(code string, msg string, err ...error) Error {
	return NewError(http.StatusBadRequest, code, msg, err...)
}

func NotFoundError(code string, msg string, err ...error) Error {
	return NewError(http.StatusBadRequest, code, msg, err...)
}

func InternalServerError(code string, msg string, err ...error) Error {
	return NewError(http.StatusInternalServerError, code, msg, err...)
}

func UnauthorizedError(code string, msg string, err ...error) Error {
	return NewError(http.StatusUnauthorized, code, msg, err...)
}

func ForbiddenError(code string, msg string, err ...error) Error {
	return NewError(http.StatusForbidden, code, msg, err...)
}

func ConflictError(code string, msg string, err ...error) Error {
	return NewError(http.StatusConflict, code, msg, err...)
}
