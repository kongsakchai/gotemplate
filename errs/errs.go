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
	Err error
	At  string
}

func (e *Error) RawError() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *Error) AtError() string {
	return e.At
}

func (e Error) Error() string {
	if e.Err == nil {
		return e.At
	}

	return fmt.Sprintf("error: %s at %s", e.Err.Error(), e.At)
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
		Err: err,
		At:  caller(maxStackDepth),
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
