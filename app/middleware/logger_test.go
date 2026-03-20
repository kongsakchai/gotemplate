package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLogger(w io.Writer) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(w, nil))
	return logger
}

type mockReaderError struct {
	err error
}

func (m *mockReaderError) Read(p []byte) (n int, err error) {
	return 0, m.err
}

func TestLoggerMiddleware(t *testing.T) {
	t.Run("should go next when disable logger", func(t *testing.T) {
		// arrange
		next := func(ctx *echo.Context) error {
			require.NotNil(t, ctx)
			return nil
		}
		enable := false

		ctx := echotest.ContextConfig{}.ToContext(t)

		// act
		middleware := Logger(enable)
		err := middleware(next)(ctx)

		// assert
		assert.NoError(t, err)
	})

	t.Run("should log request with header type when header isn't json", func(t *testing.T) {
		// arrange
		var buf bytes.Buffer
		logger := setupLogger(&buf)

		next := func(ctx *echo.Context) error {
			require.NotNil(t, ctx)
			return nil
		}
		enable := true

		ctx := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{"custom/type"},
			},
		}.ToContext(t)

		ctx.SetLogger(logger)

		// act
		middleware := Logger(enable)
		err := middleware(next)(ctx)

		// assert
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "body=custom/type")
	})

	t.Run("should log request with JSON body when header is JSON", func(t *testing.T) {
		// arrange
		var buf bytes.Buffer
		logger := setupLogger(&buf)

		next := func(ctx *echo.Context) error {
			require.NotNil(t, ctx)
			return nil
		}
		enable := true

		ctx := echotest.ContextConfig{
			JSONBody: []byte(`{"message":"hello"}`),
		}.ToContext(t)

		ctx.SetLogger(logger)

		// act
		middleware := Logger(enable)
		err := middleware(next)(ctx)

		// assert
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "body=\"{\\\"message\\\":\\\"hello\\\"}\"")
	})

	t.Run("should log error when fail to read body", func(t *testing.T) {
		// arrange
		var buf bytes.Buffer
		logger := setupLogger(&buf)

		next := func(ctx *echo.Context) error {
			require.NotNil(t, ctx)
			return nil
		}
		enable := true

		ctx := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
			},
			Request: httptest.NewRequest(http.MethodPost, "/test", &mockReaderError{err: io.ErrUnexpectedEOF}),
		}.ToContext(t)

		ctx.SetLogger(logger)

		// act
		middleware := Logger(enable)
		err := middleware(next)(ctx)

		// assert
		assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
		assert.Contains(t, buf.String(), "failed to read request body")
	})
}
