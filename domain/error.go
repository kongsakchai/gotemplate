package domain

import "fmt"

type Error struct {
	Code    string
	Message string
	Cause   string
}

func (e Error) Error() string {
	if e.Cause == "" {
		return fmt.Sprintf("(%s) %s", e.Code, e.Message)
	}
	return fmt.Sprintf("(%s) %s: %s", e.Code, e.Message, e.Cause)
}

func (e Error) Causef(log string, args ...any) Error {
	return Error{
		Code:    e.Code,
		Message: e.Message,
		Cause:   fmt.Sprintf(log, args...),
	}
}

func NewErrorCause(code, message, log string, args ...any) Error {
	return Error{
		Code:    code,
		Message: message,
		Cause:   fmt.Sprintf(log, args...),
	}
}

func NewError(code, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
