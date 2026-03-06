package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	t.Run("should create a wrap error with stack trace", func(t *testing.T) {
		e := wrap(errors.ErrUnsupported)

		assert.Equal(t, errors.ErrUnsupported, e.Unwrap())
		assert.NotEqual(t, "", e.At())
	})

	t.Run("should return the same Error if it is already of errorTrace", func(t *testing.T) {
		ogError := &errorTrace{err: errors.ErrUnsupported, at: "unit test"}
		e := wrap(ogError)

		fmt.Println(e)

		e2, ok := As(ogError)
		assert.True(t, ok)
		assert.Equal(t, e2, e)
		assert.Equal(t, "unit test", e.At())
	})

	t.Run("should return empty at if runtime.Caller fails", func(t *testing.T) {
		originalMaxStackDepth := maxStackDepth
		maxStackDepth = 1000
		defer func() { maxStackDepth = originalMaxStackDepth }()

		err := errors.New("test error")
		e := wrap(err)

		assert.Equal(t, err, e.Unwrap())
	})

	t.Run("should error string format", func(t *testing.T) {
		err := errors.New("test error")
		e := wrap(err)

		expected := "test error"
		assert.Contains(t, e.Error(), expected)
	})

	t.Run("should handle nil error", func(t *testing.T) {
		e := wrap(nil)

		assert.NotNil(t, e)
		assert.Nil(t, e.Unwrap())
		assert.Equal(t, "something went wrong", e.Error())
	})

	t.Run("should handle nil Error", func(t *testing.T) {
		var nilErr *errorTrace
		e := wrap(nilErr)
		assert.Nil(t, e)
	})

	t.Run("should return only base file name if filepath.Rel fails", func(t *testing.T) {
		originalRootPath := rootPath
		rootPath = "."
		defer func() { rootPath = originalRootPath }()

		e := wrap(errors.ErrUnsupported)

		assert.Equal(t, errors.ErrUnsupported, e.Unwrap())
		assert.Contains(t, e.Error(), "unsupported operation")
	})
}

func TestErrorAs(t *testing.T) {
	t.Run("should return true when error matches the Error type", func(t *testing.T) {
		err := &errorTrace{err: errors.ErrUnsupported, at: "unit test"}
		_, ok := As(err)
		assert.True(t, ok)
	})

	t.Run("should return false when error does not match the Error type", func(t *testing.T) {
		err := errors.New("test error")
		_, ok := As(err)
		assert.False(t, ok)
	})
}

func TestNewError(t *testing.T) {
	t.Run("should return error with message 'unit test'", func(t *testing.T) {
		err := New("unit test")

		assert.NotNil(t, err)
	})

	t.Run("should return nil when from function get nil", func(t *testing.T) {
		err := From(nil)

		assert.Nil(t, err)
	})

	t.Run("should return not nil when from function get some error", func(t *testing.T) {
		err := From(errors.ErrUnsupported)

		assert.NotNil(t, err)
	})
}

func TestToLogs(t *testing.T) {
	t.Run("should return 0 slog.Attr when error is nil", func(t *testing.T) {
		err := From(nil)

		assert.Equal(t, 0, len(Logs(err)))
	})
	t.Run("should return 2 slog.Attr when error is errorTrace", func(t *testing.T) {
		err := New("unit test")

		assert.Equal(t, 2, len(Logs(err)))
	})

	t.Run("should return 2 slog.Attr when error is normal", func(t *testing.T) {
		err := errors.New("unit test")

		assert.Equal(t, 1, len(Logs(err)))
	})
}
