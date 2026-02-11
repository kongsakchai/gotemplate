package errs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	maxStackDepth = 3
	rootPath      = ""
)

func init() {
	if wd, err := os.Getwd(); err == nil {
		rootPath = wd
	}
}

type Error struct {
	err error
	at  string
}

func (e *Error) UnwrapError() string {
	if e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *Error) At() string {
	return e.at
}

func (e Error) Error() string {
	if e.err == nil {
		return fmt.Sprintf("error: something wrong at %s", e.at)
	}

	return fmt.Sprintf("error: %s at %s", e.err.Error(), e.at)
}

func (e *Error) Unwrap() error {
	return e.err
}

func As(err error) (*Error, bool) {
	var errType *Error
	if errors.As(err, &errType) {
		return errType, true
	}

	return nil, false
}

func New(err error) error {
	if err == nil {
		return nil
	}

	return wrap(err)
}

func wrap(err error) *Error {
	var errType *Error
	if errors.As(err, &errType) {
		return errType
	}

	return &Error{
		err: err,
		at:  caller(maxStackDepth),
	}
}

func caller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	f, err := filepath.Rel(rootPath, file)
	if err != nil {
		f = filepath.Base(file)
	}
	fn := runtime.FuncForPC(pc)

	return fmt.Sprintf("(%s:%v) %s", f, line, fn.Name())
}
