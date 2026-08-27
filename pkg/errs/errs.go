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
	err   error
	at    string
	attrs []slog.Attr
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

func (e *errorTrace) LogAttrs() []slog.Attr {
	return e.attrs
}

func (e *errorTrace) With(attrs ...slog.Attr) *errorTrace {
	e.attrs = append(e.attrs, attrs...)
	return e
}

func (e *errorTrace) WithNew(attrs ...slog.Attr) *errorTrace {
	return wrap(e.err).With(attrs...)
}

func As(err error) (*errorTrace, bool) {
	if e, ok := err.(*errorTrace); ok {
		return e, true
	}

	return nil, false
}

func New(str string, args ...any) *errorTrace {
	return wrap(fmt.Errorf(str, args...))
}

func From(err error) *errorTrace {
	if err == nil {
		return nil
	}
	return wrap(err)
}

func wrap(err error) *errorTrace {
	if errType, ok := errors.AsType[*errorTrace](err); ok {
		return errType
	}
	at := caller(maxStackDepth)
	return &errorTrace{
		err:   err,
		at:    at,
		attrs: []slog.Attr{slog.String("err", err.Error()), slog.String("at", at)},
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

	return fmt.Sprintf("(%s:%d) %s", f, line, filepath.Base(fn.Name()))
}
