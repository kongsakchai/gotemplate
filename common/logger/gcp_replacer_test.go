package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceAttr(t *testing.T) {
	testcases := []struct {
		title string
		attr  slog.Attr
		want  slog.Attr
	}{
		{"replace message key", slog.String(slog.MessageKey, "test"), slog.String("message", "test")},
		{"keep other key", slog.String("other", "value"), slog.String("other", "value")},
		{"replace level key", slog.String("level", "info"), slog.String("severity", "info")},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			got, _ := GCPKeyReplacer(nil, tc.attr)
			assert.Equal(t, tc.want, got)
		})
	}
}
