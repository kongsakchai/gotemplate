package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

var maxStackDepth = 3

type Error struct {
	Err error
	at  string
}

func (e Error) Error() string {
	return fmt.Sprintf("error: %s at %s", e.Err.Error(), e.at)
}

func New(err error) *Error {
	var errType *Error
	if errors.As(err, &errType) {
		return errType
	}

	return &Error{
		Err: err,
		at:  caller(2),
	}
}

func caller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	f := filepath.Base(file)
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d (%s)", f, line, fn.Name())
}
