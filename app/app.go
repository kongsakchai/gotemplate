package app

import (
	"net/http"
)

type Response struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type RequestContext interface {
	Bind(i any) error
	Validate(i any) error
}

func Request(ctx RequestContext, target any) error {
	if err := ctx.Bind(target); err != nil {
		return BadRequest(BadRequestCode, BadRequestMsg, err)
	}
	if err := ctx.Validate(target); err != nil {
		return BadRequest(InValidCode, InValidMsg, err, err)
	}
	return nil
}

type Context interface {
	JSON(code int, i any) (err error)
}

func Ok(ctx Context, data any, msg ...string) error {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	return ctx.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(ctx Context, data any, msg ...string) error {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	return ctx.JSON(http.StatusCreated, Response{
		Code:    SuccessCode,
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Fail(ctx Context, err Error) error {
	return ctx.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Success: false,
		Data:    err.Data,
		Message: err.Message,
	})
}
