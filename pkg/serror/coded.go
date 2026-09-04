package serror

import (
	"errors"
)

type coded struct {
	Code string
	Msg  string
}

func NewCoded(code string, msg string) *coded {
	return &coded{
		Code: code,
		Msg:  msg,
	}
}

func (c *coded) Err() error {
	return wrap(c.Code, errors.New(c.Msg))
}
