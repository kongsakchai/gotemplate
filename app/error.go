package app

import "net/http"

type Error struct {
	StatusCd int
	Code     string
	Message  string
	Error    error
}

func (e Error) IsEmpty() bool {
	return e == Error{}
}

func InternalServer(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusInternalServerError,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func BadRequest(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusBadRequest,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func NotFound(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusOK,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func Unauthorized(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusUnauthorized,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func Forbidden(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusForbidden,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func Conflict(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusConflict,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}
