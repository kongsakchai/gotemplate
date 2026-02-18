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

type Errs struct {
	err error
	at  string
}

func (e *Errs) UnwrapError() string {
	if e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *Errs) At() string {
	return e.at
}

func (e Errs) Error() string {
	if e.err == nil {
		return fmt.Sprintf("error: something wrong at %s", e.at)
	}

	return fmt.Sprintf("error: %s at %s", e.err.Error(), e.at)
}

func (e *Errs) Unwrap() error {
	return e.err
}

func As(err error) (*Errs, bool) {
	var errType *Errs
	if errors.As(err, &errType) {
		return errType, true
	}

	return nil, false
}

func New(str string) error {
	return wrap(errors.New(str))
}

func Wrap(err error) error {
	if err == nil {
		return err
	}
	return wrap(err)
}

func wrap(err error) *Errs {
	var errType *Errs
	if errors.As(err, &errType) {
		return errType
	}

	return &Errs{
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
