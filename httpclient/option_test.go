package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/stretchr/testify/assert"
)

func TestTraceOption(t *testing.T) {
	t.Run("should return header with ref key", func(t *testing.T) {
		key := "ref"
		value := "some value"

		c := New(config.Config{}, TraceOption(key))
		req, err := newRequest(context.WithValue(context.Background(), key, value), c, http.MethodGet, "", nil)

		assert.NoError(t, err)
		assert.Equal(t, req.Header.Get(key), value)
	})
}
