package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	t.Run("should create a new Error with stack trace", func(t *testing.T) {
		err := errors.New("test error")
		e := New(err)

		assert.Equal(t, "test error", e.Err.Error())
		assert.NotEqual(t, "", e.at)
	})

	t.Run("should return the same Error if it is already of type Error", func(t *testing.T) {
		originalErr := New(errors.New("original error"))
		e := New(originalErr)

		assert.Equal(t, originalErr, e)
		assert.Equal(t, "original error", e.Err.Error())
		assert.NotEqual(t, "", e.at)
	})

	t.Run("should return empty at if runtime.Caller fails", func(t *testing.T) {
		originalMaxStackDepth := maxStackDepth
		maxStackDepth = 1000
		defer func() { maxStackDepth = originalMaxStackDepth }()

		err := errors.New("test error")
		e := New(err)

		assert.Equal(t, "test error", e.Err.Error())
		assert.Equal(t, "", e.at)
	})

	t.Run("should error string format", func(t *testing.T) {
		err := errors.New("test error")
		e := New(err)

		expected := "error: test error at "
		assert.Contains(t, e.Error(), expected)
		assert.Contains(t, e.Error(), e.at) // Ensure the stack trace is included
	})
}
