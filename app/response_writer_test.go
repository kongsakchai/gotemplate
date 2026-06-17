package app

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestResponseWriter(status int) *echoResponseWriter {
	rec := httptest.NewRecorder()
	return &echoResponseWriter{
		ResponseWriter: rec,
		logger:         slog.Default(),
		ctx:            nil,
		status:         0,
		url:            "/test",
		now:            time.Now(),
	}
}

func TestEchoResponseWriter_WriteHeader(t *testing.T) {
	t.Run("should set status code", func(t *testing.T) {
		w := newTestResponseWriter(0)
		w.WriteHeader(http.StatusOK)
		assert.Equal(t, http.StatusOK, w.status)
	})

	t.Run("should set not found status code", func(t *testing.T) {
		w := newTestResponseWriter(0)
		w.WriteHeader(http.StatusNotFound)
		assert.Equal(t, http.StatusNotFound, w.status)
	})
}

func TestEchoResponseWriter_Write(t *testing.T) {
	t.Run("should write body successfully", func(t *testing.T) {
		w := newTestResponseWriter(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		n, err := w.Write([]byte("test body"))
		assert.NoError(t, err)
		assert.Equal(t, 9, n)
	})

	t.Run("should write body with error status", func(t *testing.T) {
		w := newTestResponseWriter(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		n, err := w.Write([]byte("error body"))
		assert.NoError(t, err)
		assert.Equal(t, 10, n)
	})
}
