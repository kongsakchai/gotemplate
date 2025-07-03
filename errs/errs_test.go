package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("should create a new Error with stack trace", func(t *testing.T) {
		err := errors.New("test error")
		e := wrap(err)

		assert.Equal(t, "test error", e.RawError())
		assert.NotEqual(t, "", e.AtError())
	})

	t.Run("should return the same Error if it is already of type Error", func(t *testing.T) {
		originalErr := New(errors.New("original error"))
		e := wrap(originalErr)

		assert.Equal(t, originalErr, e)
		assert.Equal(t, "original error", e.RawError())
		assert.NotEqual(t, "", e.AtError())
	})

	t.Run("should return empty at if runtime.Caller fails", func(t *testing.T) {
		originalMaxStackDepth := maxStackDepth
		maxStackDepth = 1000
		defer func() { maxStackDepth = originalMaxStackDepth }()

		err := errors.New("test error")
		e := wrap(err)

		assert.Equal(t, "test error", e.RawError())
		assert.Equal(t, "", e.AtError())
	})

	t.Run("should error string format", func(t *testing.T) {
		err := errors.New("test error")
		e := wrap(err)

		expected := "error: test error at "
		assert.Contains(t, e.Error(), expected)
		assert.Contains(t, e.Error(), e.AtError())
	})

	t.Run("should handle nil error", func(t *testing.T) {
		e := wrap(nil)

		assert.NotNil(t, e)
		assert.Equal(t, "", e.RawError())
		assert.Contains(t, e.Error(), "tRunner")
		assert.Contains(t, e.Error(), "testing.go")
	})

	t.Run("should handle nil Error", func(t *testing.T) {
		var nilErr *Error
		e := wrap(nilErr)
		assert.Nil(t, e)
	})

	t.Run("should return only base file name if filepath.Rel fails", func(t *testing.T) {
		originalRootPath := rootPath
		rootPath = "."
		defer func() { rootPath = originalRootPath }()

		err := errors.New("test error")
		e := wrap(err)

		assert.Equal(t, "test error", e.RawError())
		assert.NotEqual(t, "", e.AtError())
		assert.Contains(t, e.AtError(), "testing.go")
	})
}

func TestNewError(t *testing.T) {
	t.Run("should return nil when error is nil", func(t *testing.T) {
		err := New(nil)
		assert.Nil(t, err)
	})

	t.Run("should wrap the error with stack trace", func(t *testing.T) {
		err := errors.New("test error")
		e := New(err)

		assert.NotNil(t, e)
	})
}

func TestErrorAs(t *testing.T) {
	t.Run("should return true when error matches the Error type", func(t *testing.T) {
		err := New(errors.New("test error"))
		_, ok := As(err)
		assert.True(t, ok)
	})

	t.Run("should return false when error does not match the Error type", func(t *testing.T) {
		err := errors.New("test error")
		_, ok := As(err)
		assert.False(t, ok)
	})
}
