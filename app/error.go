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

func InternalServerError(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusInternalServerError,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func BadRequestError(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusBadRequest,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func NotFoundError(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusOK,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func UnauthorizedError(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusUnauthorized,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func ForbiddenError(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusForbidden,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}

func ConflictError(code string, msg string, err error) Error {
	return Error{
		StatusCd: http.StatusConflict,
		Code:     code,
		Message:  msg,
		Error:    err,
	}
}
