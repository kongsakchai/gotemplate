package errs

import (
	"errors"
	"fmt"
	"log/slog"
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

type errorTrace struct {
	err error
	at  string
}

// Error returns error message
func (e errorTrace) Error() string {
	if e.err == nil {
		return fmt.Sprintf("something went wrong, at: %s", e.at)
	}
	return fmt.Sprintf("error: %v, at: %s", e.err, e.at)
}

// Unwrap returns the wrapped error
func (e *errorTrace) Unwrap() error {
	return e.err
}

// Unwrap returns error line
func (e *errorTrace) At() string {
	return e.at
}

func As(err error) (*errorTrace, bool) {
	if e, ok := err.(*errorTrace); ok {
		return e, true
	}

	return nil, false
}

func New(str string, args ...any) error {
	return wrap(fmt.Errorf(str, args...))
}

func From(err error) error {
	if err == nil {
		return nil
	}
	return wrap(err)
}

func wrap(err error) *errorTrace {
	var errType *errorTrace
	if errors.As(err, &errType) {
		return errType
	}

	return &errorTrace{
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

func SlogAttr(err error) []slog.Attr {
	if err == nil {
		return []slog.Attr{}
	}

	if errType, ok := As(err); ok {
		return []slog.Attr{
			slog.String("err", errType.err.Error()),
			slog.String("at", errType.at),
		}
	}

	return []slog.Attr{
		slog.String("err", err.Error()),
	}
}
