package logger

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type replaceAttrFunc func(groups []string, a slog.Attr) slog.Attr

type handlerOptions struct {
	level       slog.Level
	replaceAttr replaceAttrFunc
	timeFormat  string
}

type textHandler struct {
	opts handlerOptions

	mu sync.Mutex
	w  io.Writer
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.level.Level()
}

func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &textHandler{}
}

func (h *textHandler) WithGroup(name string) slog.Handler {
	return &textHandler{}
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := newBuffer()
	defer buf.Free()

	buf.WriteByte('\n')
	// Write Emoji
	switch r.Level {
	case slog.LevelError:
		buf.WriteString("🚨 ")
	case slog.LevelWarn:
		buf.WriteString("⚠️  ")
	case slog.LevelInfo:
		buf.WriteString("🌱 ")
	case slog.LevelDebug:
		buf.WriteString("🐛 ")
	}

	// Write Time
	if !r.Time.IsZero() {
		val := r.Time.Round(0)
		buf.WriteString("\u001b[38;2;165;173;203m")
		*buf = val.AppendFormat(*buf, h.opts.timeFormat)
		buf.WriteString("\u001b[0m ")
	}

	// Write Level
	switch r.Level {
	case slog.LevelError:
		buf.WriteString("\u001b[38;2;237;135;160mERROR\u001b[0m ")
	case slog.LevelWarn:
		buf.WriteString("\u001b[38;2;238;212;159mWARN\u001b[0m ")
	case slog.LevelInfo:
		buf.WriteString("\u001b[38;2;166;218;149mINFO\u001b[0m ")
	case slog.LevelDebug:
		buf.WriteString("DEBUG ")
	}

	// Write Message
	if r.Message != "" {
		buf.WriteString("\u001b[1m")
		buf.WriteString(r.Message)
		buf.WriteString("\u001b[0m ")
	}

	// Write Attributes
	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			h.appendAttr(buf, a)
			return true
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*buf)
	return err
}

func (h *textHandler) appendAttr(buf *buffer, a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}

	if a.Value.Kind() == slog.KindGroup {
	} else {
		h.appendKey(buf, a.Key)
		h.appendValue(buf, a.Value)
	}

	buf.WriteByte(' ')
}

func (h *textHandler) appendKey(buf *buffer, key string) {
	buf.WriteString(key)
	buf.WriteByte('=')
}

func (h *textHandler) appendValue(buf *buffer, v slog.Value) {
	switch v.Kind() {
	case slog.KindString:
		buf.WriteString(v.String())
	case slog.KindTime:
		*buf = appendRFC3339Millis(*buf, v.Time())
	case slog.KindBool:
		*buf = strconv.AppendBool(*buf, v.Bool())
	case slog.KindInt64:
		*buf = strconv.AppendInt(*buf, v.Int64(), 10)
	case slog.KindUint64:
		*buf = strconv.AppendUint(*buf, v.Uint64(), 10)
	case slog.KindFloat64:
		*buf = strconv.AppendFloat(*buf, v.Float64(), 'g', -1, 64)
	case slog.KindDuration:
		*buf = strconv.AppendInt(*buf, int64(v.Duration()), 10)
	case slog.KindAny:
		defer func() {
			// Copied from log/slog/handler.go.
			if r := recover(); r != nil {
				// If it panics with a nil pointer, the most likely cases are
				// an encoding.TextMarshaler or error fails to guard against nil,
				// in which case "<nil>" seems to be the feasible choice.
				//
				// Adapted from the code in fmt/print.go.
				if v := reflect.ValueOf(v.Any()); v.Kind() == reflect.Pointer && v.IsNil() {
					buf.WriteString("<nil>")
					return
				}

				// Otherwise just print the original panic message.
				buf.WriteString(fmt.Sprintf("<panic: %v>", r))
			}
		}()

		switch cv := v.Any().(type) {
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				break
			}
			buf.Write(data)
		default:
			buf.WriteString(fmt.Sprintf("%+v", v.Any()))
		}
	}
}

// copy from log/slog/handler.go
func appendRFC3339Millis(b []byte, t time.Time) []byte {
	// Format according to time.RFC3339Nano since it is highly optimized,
	// but truncate it to use millisecond resolution.
	// Unfortunately, that format trims trailing 0s, so add 1/10 millisecond
	// to guarantee that there are exactly 4 digits after the period.
	const prefixLen = len("2006-01-02T15:04:05.000")
	n := len(b)
	t = t.Truncate(time.Millisecond).Add(time.Millisecond / 10)
	b = t.AppendFormat(b, time.RFC3339Nano)
	b = append(b[:n+prefixLen], b[n+prefixLen+1:]...) // drop the 4th digit
	return b
}
