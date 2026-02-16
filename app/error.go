package app

import (
	"fmt"
	"net/http"
)

type Error struct {
	HTTPCode int
	Code     string
	Message  string
	Data     any
	Err      error
}

func (e Error) IsEmpty() bool {
	return e == Error{}
}

func (e Error) Error() string {
	return fmt.Sprintf("http_code=%d code=%s msg=\"%s\" data=%v err=%v", e.HTTPCode, e.Code, e.Message, e.Data, e.Err)
}

func errorData(data []any) any {
	if len(data) == 0 {
		return nil
	}
	if len(data) == 1 {
		return data[0]
	}
	return data
}

func InternalServer(code string, msg string, err error, data ...any) Error {
	return Error{
		HTTPCode: http.StatusInternalServerError,
		Code:     code,
		Message:  msg,
		Err:      err,
		Data:     errorData(data),
	}
}

func BadRequest(code string, msg string, err error, data ...any) Error {
	return Error{
		HTTPCode: http.StatusBadRequest,
		Code:     code,
		Message:  msg,
		Err:      err,
		Data:     errorData(data),
	}
}

func NotFound(code string, msg string, err error, data ...any) Error {
	return Error{
		HTTPCode: http.StatusOK,
		Code:     code,
		Message:  msg,
		Err:      err,
		Data:     errorData(data),
	}
}

func Unauthorized(code string, msg string, err error, data ...any) Error {
	return Error{
		HTTPCode: http.StatusUnauthorized,
		Code:     code,
		Message:  msg,
		Err:      err,
		Data:     errorData(data),
	}
}

func Forbidden(code string, msg string, err error, data ...any) Error {
	return Error{
		HTTPCode: http.StatusForbidden,
		Code:     code,
		Message:  msg,
		Err:      err,
		Data:     errorData(data),
	}
}

func Conflict(code string, msg string, err error, data ...any) Error {
	return Error{
		HTTPCode: http.StatusConflict,
		Code:     code,
		Message:  msg,
		Err:      err,
		Data:     errorData(data),
	}
}
