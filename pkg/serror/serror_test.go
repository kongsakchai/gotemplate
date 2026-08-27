package serror

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("should set rootPath from os.Getwd", func(t *testing.T) {
		assert.NotEmpty(t, rootPath)
	})
}

func TestErrorMethod(t *testing.T) {
	t.Run("should return error string formatted", func(t *testing.T) {
		err := errors.New("test error")
		e := From(err)

		expected := "test error"
		assert.Contains(t, e.Error(), expected)
	})

	t.Run("should contain at field in error string", func(t *testing.T) {
		err := errors.New("test error")
		e := From(err)

		assert.Contains(t, e.Error(), ", at: ")
	})
}

func TestMessageMethod(t *testing.T) {
	t.Run("should return wrapped error message", func(t *testing.T) {
		err := errors.New("wrapped message")
		e := From(err)

		assert.Equal(t, "wrapped message", e.Message())
	})
}

func TestUnwrapMethod(t *testing.T) {
	t.Run("should return wrapped error", func(t *testing.T) {
		base := errors.New("base error")
		e := From(base)

		assert.Equal(t, base, e.Unwrap())
	})
}

func TestCodeAtMethod(t *testing.T) {
	t.Run("should return code set by NewCode", func(t *testing.T) {
		e := NewWithCode("E001", "error message")

		assert.Equal(t, "E001", e.Code())
	})

	t.Run("should return empty code set by New", func(t *testing.T) {
		e := New("error message")

		assert.Equal(t, "", e.Code())
	})

	t.Run("should return empty code set by From", func(t *testing.T) {
		e := From(errors.New("error"))

		assert.Equal(t, "", e.Code())
	})
}

func TestAtMethod(t *testing.T) {
	t.Run("should return non-empty stack trace location", func(t *testing.T) {
		e := From(errors.ErrUnsupported)

		assert.NotEqual(t, "", e.At())
	})

	t.Run("should return custom at value", func(t *testing.T) {
		e := &serror{err: errors.New("test"), at: "custom location"}

		assert.Equal(t, "custom location", e.At())
	})
}

func TestIsMethod(t *testing.T) {
	t.Run("should return true when wrapped error matches target using As", func(t *testing.T) {
		base := errors.ErrUnsupported
		e := From(base)

		assert.True(t, e.Is(base))
	})

	t.Run("should return true when wrapped error matches sentinel error", func(t *testing.T) {
		e := From(errors.ErrUnsupported)

		assert.True(t, e.Is(errors.ErrUnsupported))
	})

	t.Run("should return false when wrapped error does not match", func(t *testing.T) {
		e := From(errors.New("different error"))

		assert.False(t, e.Is(errors.ErrUnsupported))
	})

	t.Run("should return true when Is checks against a *serror target", func(t *testing.T) {
		e := From(errors.ErrUnsupported)
		target := &serror{err: errors.ErrUnsupported, at: "target"}

		assert.True(t, e.Is(target))
	})
}

func TestLogAttrsMethod(t *testing.T) {
	t.Run("should return log attributes initialized by wrap", func(t *testing.T) {
		err := errors.New("test error")
		e := From(err)

		attrs := e.LogAttrs()
		assert.NotEmpty(t, attrs)

		var errAttrFound, atAttrFound bool
		for _, attr := range attrs {
			if attr.Key == "err" {
				assert.Equal(t, "test error", attr.Value.String())
				errAttrFound = true
			}
			if attr.Key == "at" {
				assert.NotEmpty(t, attr.Value.String())
				atAttrFound = true
			}
		}
		assert.True(t, errAttrFound, "expected 'err' attribute")
		assert.True(t, atAttrFound, "expected 'at' attribute")
	})
}

func TestWithMethod(t *testing.T) {
	t.Run("should append additional attributes and return self", func(t *testing.T) {
		e := From(errors.New("test"))

		newErr := e.With(slog.String("key", "value"))
		assert.Same(t, e, newErr)

		attrs := e.LogAttrs()
		assert.Len(t, attrs, 3) // err + at + key

		var foundKey bool
		for _, attr := range attrs {
			if attr.Key == "key" {
				assert.Equal(t, "value", attr.Value.String())
				foundKey = true
			}
		}
		assert.True(t, foundKey, "expected 'key' attribute after With")
	})

	t.Run("should chain multiple With calls", func(t *testing.T) {
		e := From(errors.New("test")).
			With(slog.String("a", "1")).
			With(slog.Int("b", 2))

		attrs := e.LogAttrs()
		assert.Len(t, attrs, 4) // err + at + a + b
	})

	t.Run("should work with With on zero-value-like error", func(t *testing.T) {
		e := &serror{err: errors.New("minimal")}
		e = e.With(slog.String("x", "y"))

		attrs := e.LogAttrs()
		assert.Len(t, attrs, 1) // x
	})
}

func TestErrorAs(t *testing.T) {
	t.Run("should return true when error matches the Error type", func(t *testing.T) {
		err := &serror{err: errors.ErrUnsupported, at: "unit test"}
		_, ok := As(err)
		assert.True(t, ok)
	})

	t.Run("should return false when error does not match the Error type", func(t *testing.T) {
		err := errors.New("test error")
		_, ok := As(err)
		assert.False(t, ok)
	})
}

func TestNewFunction(t *testing.T) {
	t.Run("should return error with message", func(t *testing.T) {
		err := New("unit test")

		assert.NotNil(t, err)
		assert.Equal(t, "unit test", err.Message())
		assert.Equal(t, "", err.Code())
	})

	t.Run("should support format arguments", func(t *testing.T) {
		err := New("user %s not found", "alice")

		assert.NotNil(t, err)
		assert.Equal(t, "user alice not found", err.Message())
	})
}

func TestNewCodeFunction(t *testing.T) {
	t.Run("should create error with code and message", func(t *testing.T) {
		err := NewWithCode("E001", "database error")

		assert.NotNil(t, err)
		assert.Equal(t, "E001", err.Code())
		assert.Equal(t, "database error", err.Message())
	})

	t.Run("should support format arguments with code", func(t *testing.T) {
		err := NewWithCode("E002", "user %s not found", "bob")

		assert.NotNil(t, err)
		assert.Equal(t, "E002", err.Code())
		assert.Equal(t, "user bob not found", err.Message())
	})
}

func TestFromFunction(t *testing.T) {
	t.Run("should return nil when given nil", func(t *testing.T) {
		assert.Nil(t, From(nil))
	})

	t.Run("should wrap a non-nil error", func(t *testing.T) {
		err := From(errors.ErrUnsupported)

		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrUnsupported, err.Unwrap())
		assert.Equal(t, "", err.Code())
	})

	t.Run("should handle nil *serror", func(t *testing.T) {
		var nilErr *serror
		assert.Nil(t, From(nilErr))
	})
}

func TestFromCodeFunction(t *testing.T) {
	t.Run("should create error with code from existing error", func(t *testing.T) {
		err := FromWithCode("E003", errors.New("from code error"))

		assert.NotNil(t, err)
		assert.Equal(t, "E003", err.Code())
		assert.Equal(t, "from code error", err.Message())
	})

	t.Run("should wrap ErrUnsupported with code", func(t *testing.T) {
		err := FromWithCode("E004", errors.ErrUnsupported)

		assert.NotNil(t, err)
		assert.Equal(t, "E004", err.Code())
		assert.Equal(t, errors.ErrUnsupported, err.Unwrap())
	})
}

func TestWrapFunction(t *testing.T) {
	t.Run("should return nil for nil error", func(t *testing.T) {
		result := wrap("", nil)
		assert.Nil(t, result)
	})

	t.Run("should create serror with provided code", func(t *testing.T) {
		result := wrap("MYCODE", errors.New("test"))

		assert.NotNil(t, result)
		assert.Equal(t, "MYCODE", result.Code())
		assert.Equal(t, "test", result.Message())
		assert.NotEmpty(t, result.At())
	})

	t.Run("should return error with proper formatting", func(t *testing.T) {
		result := wrap("", errors.New("wrap test"))

		assert.Contains(t, result.Error(), "wrap test")
		assert.Contains(t, result.Error(), ", at: ")
	})

	t.Run("should include err and at in LogAttrs", func(t *testing.T) {
		result := wrap("code123", errors.New("attr test"))

		attrs := result.LogAttrs()
		assert.GreaterOrEqual(t, len(attrs), 2)
	})
}

func TestCallerFunction(t *testing.T) {
	t.Run("should return valid caller info string", func(t *testing.T) {
		info := caller(maxStackDepth)
		assert.NotEmpty(t, info)
		// Should contain parentheses with file info
		assert.Contains(t, info, "(")
		assert.Contains(t, info, ")")
	})

	t.Run("should return empty string when skip exceeds call depth", func(t *testing.T) {
		originalMaxStackDepth := maxStackDepth
		maxStackDepth = 1000
		defer func() { maxStackDepth = originalMaxStackDepth }()

		info := caller(1000)
		assert.Equal(t, "", info)
	})

	t.Run("should use filepath.Base when filepath.Rel fails", func(t *testing.T) {
		originalRootPath := rootPath
		rootPath = "." // unlikely to match actual file paths, may trigger fallback
		defer func() { rootPath = originalRootPath }()

		info := caller(maxStackDepth)
		assert.NotEmpty(t, info)
	})
}
