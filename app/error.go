package app

import "net/http"

type Error struct {
	HTTPCode int
	Code     string
	Message  string
	Error    error
}

func (e Error) IsEmpty() bool {
	return e == Error{}
}

func InternalServer(code string, msg string, err error) Error {
	return Error{
		HTTPCode: http.StatusInternalServerError,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func BadRequest(code string, msg string, err error) Error {
	return Error{
		HTTPCode: http.StatusBadRequest,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func NotFound(code string, msg string, err error) Error {
	return Error{
		HTTPCode: http.StatusOK,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func Unauthorized(code string, msg string, err error) Error {
	return Error{
		HTTPCode: http.StatusUnauthorized,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func Forbidden(code string, msg string, err error) Error {
	return Error{
		HTTPCode: http.StatusForbidden,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func Conflict(code string, msg string, err error) Error {
	return Error{
		HTTPCode: http.StatusConflict,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}
