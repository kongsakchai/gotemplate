package app

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

func newTestResponseWriter(t *testing.T, status int) *echoResponseWriter {
	ctxConfig := echotest.ContextConfig{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/test",
			},
		},
	}

	ctx, rec := ctxConfig.ToContextRecorder(t)
	return &echoResponseWriter{
		ResponseWriter: rec,
		ctx:            ctx,
		status:         0,
	}
}

func TestEchoResponseWriter_WriteHeader(t *testing.T) {
	t.Run("should set status code", func(t *testing.T) {
		w := newTestResponseWriter(t, 0)
		w.WriteHeader(http.StatusOK)
		assert.Equal(t, http.StatusOK, w.status)
	})

	t.Run("should set not found status code", func(t *testing.T) {
		w := newTestResponseWriter(t, 0)
		w.WriteHeader(http.StatusNotFound)
		assert.Equal(t, http.StatusNotFound, w.status)
	})
}

func TestEchoResponseWriter_Write(t *testing.T) {
	t.Run("should write body successfully", func(t *testing.T) {
		w := newTestResponseWriter(t, http.StatusOK)
		w.WriteHeader(http.StatusOK)
		n, err := w.Write([]byte("test body"))
		assert.NoError(t, err)
		assert.Equal(t, 9, n)
	})

	t.Run("should write body with error status", func(t *testing.T) {
		w := newTestResponseWriter(t, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		n, err := w.Write([]byte("error body"))
		assert.NoError(t, err)
		assert.Equal(t, 10, n)
	})
}
