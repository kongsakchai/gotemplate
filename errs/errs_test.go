package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("should create a wrap error with stack trace", func(t *testing.T) {
		err := errors.New("test error")
		e := wrap(err)

		assert.Equal(t, "test error", e.UnwrapError())
		assert.NotEqual(t, "", e.At())
	})

	t.Run("should return the same Error if it is already of type Error", func(t *testing.T) {
		originalErr := New("original error")
		e := wrap(originalErr)

		fmt.Println(e)

		assert.Equal(t, originalErr, e)
		assert.Equal(t, "original error", e.UnwrapError())
		assert.NotEqual(t, "", e.At())
	})

	t.Run("should return empty at if runtime.Caller fails", func(t *testing.T) {
		originalMaxStackDepth := maxStackDepth
		maxStackDepth = 1000
		defer func() { maxStackDepth = originalMaxStackDepth }()

		err := errors.New("test error")
		e := wrap(err)

		assert.Equal(t, "test error", e.UnwrapError())
		assert.Equal(t, "", e.At())
	})

	t.Run("should error string format", func(t *testing.T) {
		err := errors.New("test error")
		e := wrap(err)

		expected := "error: test error at "
		assert.Contains(t, e.Error(), expected)
		assert.Contains(t, e.Error(), e.At())
	})

	t.Run("should handle nil error", func(t *testing.T) {
		e := wrap(nil)

		assert.NotNil(t, e)
		assert.Equal(t, "", e.UnwrapError())
		assert.Contains(t, e.Error(), "tRunner")
		assert.Contains(t, e.Error(), "testing.go")
	})

	t.Run("should handle nil Error", func(t *testing.T) {
		var nilErr *Errs
		e := wrap(nilErr)
		assert.Nil(t, e)
	})

	t.Run("should return only base file name if filepath.Rel fails", func(t *testing.T) {
		originalRootPath := rootPath
		rootPath = "."
		defer func() { rootPath = originalRootPath }()

		err := errors.New("test error")
		e := wrap(err)

		assert.Equal(t, "test error", e.Unwrap().Error())
		assert.NotEqual(t, "", e.At())
		assert.Contains(t, e.Error(), "testing.go")
	})
}

func TestWrapError(t *testing.T) {
	t.Run("should return nil when error is nil", func(t *testing.T) {
		err := Wrap(nil)
		assert.Nil(t, err)
	})

	t.Run("should wrap the error with stack trace", func(t *testing.T) {
		e := Wrap(errors.New("error"))

		assert.NotNil(t, e)
	})
}

func TestErrorAs(t *testing.T) {
	t.Run("should return true when error matches the Error type", func(t *testing.T) {
		err := New("test error")
		_, ok := As(err)
		assert.True(t, ok)
	})

	t.Run("should return false when error does not match the Error type", func(t *testing.T) {
		err := errors.New("test error")
		_, ok := As(err)
		assert.False(t, ok)
	})
}
