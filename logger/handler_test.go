package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestEnabled(t *testing.T) {
	h := NewTextHandler(nil, &HandlerOptions{
		Level: slog.LevelInfo,
	})

	if !h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("Expected Enabled to return true for LevelInfo")
	}

	if h.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("Expected Enabled to return false for LevelDebug")
	}
}

func TestWithAttrs(t *testing.T) {
	h := NewTextHandler(nil, nil)

	attrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
	}

	newHandler := h.WithAttrs(attrs)

	if len(newHandler.(*textHandler).groups) != 0 {
		t.Error("Expected groups to be empty in new handler")
	}

	buf := newBuffer()
	for _, a := range attrs {
		h.en.writeAttr(buf, h.groups, a)
	}

	if string(*buf) == "" {
		t.Error("Expected buffer to contain attributes")
	}
}

func TestWithGroup(t *testing.T) {
	h := NewTextHandler(nil, nil)

	groupName := "testGroup"
	newHandler := h.WithGroup(groupName)

	if len(newHandler.(*textHandler).groups) != 1 || newHandler.(*textHandler).groups[0] != groupName {
		t.Error("Expected groups to contain the group name in new handler")
	}

	newHandler = newHandler.WithGroup("anotherGroup")
	if len(newHandler.(*textHandler).groups) != 2 || newHandler.(*textHandler).groups[1] != "anotherGroup" {
		t.Error("Expected groups to contain both group names in new handler")
	}
}

func TestHandler(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	h := NewTextHandler(b, &HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: time.RFC3339,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "key1" {
				return slog.String("key1", "replaced-value")
			}
			return a
		},
	})
	h.attrPrefix = []byte("prefix=\"prefix-data\"")

	rec := slog.NewRecord(
		time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		slog.LevelInfo,
		"test message",
		uintptr(0),
	)
	rec.AddAttrs(slog.String("key1", "value1"))

	err := h.Handle(context.Background(), rec)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if b.Len() == 0 {
		t.Error("Expected buffer to contain log output")
	}

	if !bytes.Contains(b.Bytes(), []byte("test message")) {
		t.Error("Expected buffer to contain the log message")
	}

	// if !bytes.Contains(b, []byte("test message")) {
	// 	t.Error("Expected buffer to contain the log message")
	// }
}
