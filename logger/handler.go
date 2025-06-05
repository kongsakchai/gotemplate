package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type replaceAttrFunc func(groups []string, a slog.Attr) slog.Attr

type HandlerOptions struct {
	Level       slog.Level
	ReplaceAttr replaceAttrFunc
	TimeFormat  string
}

type textHandler struct {
	opts       HandlerOptions
	attrPrefix []byte
	groups     []string

	mu sync.Mutex
	w  io.Writer
}

func NewTextHandler(w io.Writer, opts *HandlerOptions) *textHandler {
	if opts == nil {
		opts = &HandlerOptions{
			Level:      slog.LevelInfo,
			TimeFormat: time.RFC3339,
		}
	}

	return &textHandler{
		opts:   *opts,
		w:      w,
		groups: make([]string, 0),
	}
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	buf := newBuffer()
	for _, a := range attrs {
		h.appendAttr(buf, a)
	}
	return &textHandler{
		opts:       h.opts,
		w:          h.w,
		attrPrefix: *buf,
		groups:     h.groups,
	}
}

func (h *textHandler) WithGroup(name string) slog.Handler {
	gs := make([]string, len(h.groups)+1)
	if len(h.groups) > 0 {
		copy(gs, h.groups)
	}
	gs[len(gs)-1] = name

	return &textHandler{
		opts:       h.opts,
		w:          h.w,
		attrPrefix: h.attrPrefix,
		groups:     gs,
	}
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := newBuffer()
	defer buf.Free()

	// Write Emoji
	switch r.Level {
	case slog.LevelError:
		buf.WriteString("❌ ")
	case slog.LevelWarn:
		buf.WriteString("⚠️  ")
	case slog.LevelInfo:
		buf.WriteString("🌱 ")
	case slog.LevelDebug:
		buf.WriteString("🐛 ")
	}

	// Write Time
	if !r.Time.IsZero() {
		var val time.Time
		if rep := h.opts.ReplaceAttr; rep != nil {
			t := rep(h.groups, slog.Time(slog.TimeKey, r.Time))
			val = t.Value.Time().Round(0)
		} else {
			val = r.Time.Round(0)
		}

		buf.WriteString(colorGray)
		*buf = val.AppendFormat(*buf, h.opts.TimeFormat)
		buf.WriteString(colorResetWithSpace)
	}

	// Write Level
	switch r.Level {
	case slog.LevelError:
		buf.WriteString(colorRed + "ERR" + colorResetWithSpace)
	case slog.LevelWarn:
		buf.WriteString(colorYellow + "WRN" + colorResetWithSpace)
	case slog.LevelInfo:
		buf.WriteString(colorGreen + "INF" + colorResetWithSpace)
	case slog.LevelDebug:
		buf.WriteString("DBG")
	}

	// Write Message
	if r.Message != "" {
		buf.WriteString(qouteText + boldText)
		buf.WriteString(r.Message)
		buf.WriteString(colorReset + qouteText + " ")
	}

	// Wrote attrPrefix
	prefix := h.attrPrefix
	if len(prefix) > 0 {
		buf.Write(h.attrPrefix)
	}

	// Write Attributes
	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			h.appendAttr(buf, a)
			return true
		})
	}

	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*buf)
	return err
}

func (h *textHandler) appendAttr(buf *buffer, a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}

	if rep := h.opts.ReplaceAttr; rep != nil && a.Value.Kind() == slog.KindGroup {
		a = rep(h.groups, a)
		a.Value = a.Value.Resolve()
	}

	if a.Value.Kind() == slog.KindGroup {
		for _, group := range a.Value.Group() {
			h.openGroup(a.Key)
			h.appendAttr(buf, group)
			h.closeGroup()
		}
	} else {
		h.appendKey(buf, a.Key)
		h.appendValue(buf, a.Value)
	}
}

func (h *textHandler) openGroup(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.groups = append(h.groups, key)
}

func (h *textHandler) closeGroup() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.groups) > 0 {
		h.groups = h.groups[:len(h.groups)-1]
	}
}

func (h *textHandler) appendKey(buf *buffer, keys string) {
	buf.WriteString(colorCyan)
	for _, key := range h.groups {
		buf.WriteString(key)
		buf.WriteByte('.')
	}

	buf.WriteString(keys)
	buf.WriteString(colorReset)
	buf.WriteByte('=')
}

func (h *textHandler) appendValue(buf *buffer, v slog.Value) {
	defer func() {
		// Copied from log/slog/handler.go.
		if r := recover(); r != nil {
			// If it panics with a nil pointer, the most likely cases are
			// an encoding.TextMarshaler or error fails to guard against nil,
			// in which case "<nil>" seems to be the feasible choice.
			//
			// Adapted from the code in fmt/print.go.
			if v := reflect.ValueOf(v.Any()); v.Kind() == reflect.Pointer && v.IsNil() {
				buf.WriteString("nil")
				return
			}

			// Otherwise just print the original panic message.
			buf.WriteString(fmt.Sprintf("<panic: %v>", r))
		}
	}()

	kind := v.Kind()
	switch kind {
	case slog.KindString:
		*buf = strconv.AppendQuote(*buf, v.String())
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
		*buf = appendWithRefect(*buf, reflect.ValueOf(v.Any()))
	}

	buf.WriteByte(' ')
}

func appendWithRefect(b []byte, v reflect.Value) []byte {
	kind := v.Kind()
	switch kind {
	case reflect.String:
		return append(b, v.String()...)
	case reflect.Bool:
		return strconv.AppendBool(b, v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(b, v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(b, v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.AppendFloat(b, v.Float(), 'g', -1, 64)
	case reflect.Slice, reflect.Array:
		if kind == reflect.Slice && v.IsNil() {
			return append(b, "nil"...)
		}

		b = append(b, '[')
		for i := range v.Len() {
			if i > 0 {
				b = append(b, ',', ' ')
			}
			b = appendWithRefect(b, v.Index(i))
		}

		return append(b, ']')
	case reflect.Map:
		if v.IsNil() {
			return append(b, "nil"...)
		}

		b = append(b, '{')
		for i, key := range v.MapKeys() {
			if i > 0 {
				b = append(b, ',', ' ')
			}
			b = appendWithRefect(b, key)
			b = append(b, ':')
			b = appendWithRefect(b, v.MapIndex(key))
		}

		return append(b, '}')
	case reflect.Struct:
		b = append(b, '{')
		for i := range v.NumField() {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}

			if i > 0 {
				b = append(b, ',', ' ')
			}

			b = append(b, field.Name...)
			b = append(b, ':')
			b = appendWithRefect(b, v.Field(i))
		}

		return append(b, '}')
	case reflect.Ptr:
		if v.IsNil() {
			return append(b, "nil"...)
		}

		return appendWithRefect(b, v.Elem())
	}

	return b
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
