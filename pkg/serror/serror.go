package serror

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

type serror struct {
	err  error
	code string
	at   string

	Attrs []slog.Attr
	Data  []any
}

// Error returns the error message
func (e serror) Error() string {
	return fmt.Sprintf("error: %v, code: %s, at: %s", e.err, e.code, e.at)
}

// Msg returns the wrapped error message
func (e *serror) Msg() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error
func (e *serror) Unwrap() error {
	return e.err
}

// Code returns the error code
func (e *serror) Code() string {
	return e.code
}

// At returns error line
func (e *serror) At() string {
	return e.at
}

// LogAttrs returns all slog attributes; implement logger.LogAttributer
func (e *serror) LogAttrs() []slog.Attr {
	return e.Attrs
}

// WithAttr appends slog attributes and returns the same error for chaining
func (e *serror) WithAttr(attrs ...slog.Attr) *serror {
	e.Attrs = append(e.Attrs, attrs...)
	return e
}

// WithData appends data fields and returns the same error for chaining
func (e *serror) WithData(data ...any) *serror {
	e.Data = append(e.Data, data...)
	return e
}

// Implement errors.Is
func (e *serror) Is(err error) bool {
	if target, ok := As(err); ok {
		return errors.Is(e.err, target.err)
	}
	return errors.Is(e.err, err)
}

func As(err error) (*serror, bool) {
	return errors.AsType[*serror](err)
}

// Create
func New(str string, args ...any) *serror {
	return wrap("", fmt.Errorf(str, args...))
}

func NewWithCode(code string, str string, args ...any) *serror {
	return wrap(code, fmt.Errorf(str, args...))
}

func From(err error) *serror {
	return wrap("", err)
}

func FromWithCode(code string, err error) *serror {
	return wrap(code, err)
}

func wrap(code string, err error) *serror {
	if err == nil {
		return nil
	}
	if serr, ok := As(err); ok && serr == nil { // handle typed nil *serror
		return nil
	}
	at := caller(maxStackDepth)
	return &serror{
		err:   err,
		code:  code,
		at:    at,
		Attrs: []slog.Attr{slog.String("err", err.Error()), slog.String("at", at)},
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
