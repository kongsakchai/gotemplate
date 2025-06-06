package logger

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"time"
)

const (
	txtGray    = "\033[90m"
	txtRed     = "\033[31m"
	txtGreen   = "\033[32m"
	txtYellow  = "\033[33m"
	txtBlue    = "\033[34m"
	txtMagenta = "\033[35m"
	txtCyan    = "\033[36m"
	txtWhite   = "\033[37m"

	txtBold  = "\033[1m"
	txtReset = "\033[0m"
)

type encodeText struct {
	opt HandlerOptions
}

func newEncodeText(opt HandlerOptions) *encodeText {
	return &encodeText{
		opt: opt,
	}
}

func (e *encodeText) style(buf *buffer, color string) {
	if e.opt.DisableColor {
		return
	}
	buf.WriteString(color)
}

func (e *encodeText) reset(buf *buffer) {
	if e.opt.DisableColor {
		return
	}
	buf.WriteString(txtReset)
}

func (e *encodeText) writeEmojiLevel(buf *buffer, level slog.Level) {
	if e.opt.DisableEmoji {
		return
	}

	switch level {
	case slog.LevelError:
		buf.WriteString("❌ ")
	case slog.LevelWarn:
		buf.WriteString("⚠️  ")
	case slog.LevelInfo:
		buf.WriteString("🌱 ")
	case slog.LevelDebug:
		buf.WriteString("🐛 ")
	}
}

func (e *encodeText) writeTime(buf *buffer, t time.Time) {
	if e.opt.DisableTime {
		return
	}

	e.style(buf, txtGray)
	*buf = t.AppendFormat(*buf, e.opt.TimeFormat)
	e.reset(buf)
	buf.WriteByte(' ')
}

func (e *encodeText) writeLevel(buf *buffer, level slog.Level) {
	if e.opt.DisableLevel {
		return
	}

	switch level {
	case slog.LevelError:
		e.style(buf, txtRed)
		buf.WriteString("ERR ")
		e.reset(buf)
	case slog.LevelWarn:
		e.style(buf, txtYellow)
		buf.WriteString("WRN ")
		e.reset(buf)
	case slog.LevelInfo:
		e.style(buf, txtGreen)
		buf.WriteString("INF ")
		e.reset(buf)
	case slog.LevelDebug:
		buf.WriteString("DBG ")
	}

}

func (e *encodeText) writeMessage(buf *buffer, str string) {
	e.style(buf, txtBold)
	e.writeString(buf, str)
	e.reset(buf)
	e.writeSpace(buf)
}

func (e *encodeText) writeAttr(buf *buffer, gs []string, a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}

	val := a.Value.Resolve()
	if val.Kind() == slog.KindGroup {
		gs = append(gs, a.Key)
		for _, subAttr := range val.Group() {
			e.writeAttr(buf, gs, subAttr)
		}
		gs = gs[:len(gs)-1]
	} else {
		e.writeKey(buf, gs, a.Key)
		e.writeValue(buf, val)
	}
}

func (e *encodeText) writeKey(buf *buffer, gs []string, key string) {
	e.style(buf, txtCyan)
	for _, g := range gs {
		buf.WriteString(g)
		buf.WriteByte('.')
	}

	buf.WriteString(key)
	e.style(buf, txtGray)
	buf.WriteByte('=')
	e.reset(buf)
}

func (e *encodeText) writeValue(buf *buffer, val slog.Value) {
	switch val.Kind() {
	case slog.KindBool:
		e.writeBool(buf, val.Bool())
	case slog.KindInt64:
		e.writeInt(buf, val.Int64())
	case slog.KindUint64:
		e.writeUint(buf, val.Uint64(), 10)
	case slog.KindFloat64:
		e.writeFloat(buf, val.Float64())
	case slog.KindString:
		e.writeString(buf, val.String())
	case slog.KindTime:
		e.writeTimeRFC3339(buf, val.Time())
	case slog.KindDuration:
		e.writeDuration(buf, val.Duration())
	case slog.KindAny:
		e.writeAny(buf, reflect.ValueOf(val.Any()))
	}
	e.writeSpace(buf)
}

func (e *encodeText) writeAny(buf *buffer, val reflect.Value) {
	switch val.Kind() {
	case reflect.Bool:
		e.writeBool(buf, val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.writeInt(buf, val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.writeUint(buf, val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		e.writeFloat(buf, val.Float())
	case reflect.String:
		e.writeString(buf, val.String())
	case reflect.Struct:
		if str, ok := val.Interface().(fmt.Stringer); ok {
			buf.WriteString(str.String())
			return
		}

		buf.WriteByte('{')
		for i := range val.NumField() {
			field := val.Type().Field(i)
			if i > 0 {
				e.writeSpace(buf)
			}
			buf.WriteString(field.Name)
			buf.WriteByte(':')
			e.writeAny(buf, val.Field(i))
		}
		buf.WriteByte('}')
	case reflect.Slice, reflect.Array:
		buf.WriteByte('[')
		for i := range val.Len() {
			if i > 0 {
				e.writeSpace(buf)
			}
			e.writeAny(buf, val.Index(i))
		}
		buf.WriteByte(']')
	case reflect.Map:
		buf.WriteByte('[')
		keys := val.MapKeys()
		for i, key := range keys {
			if i > 0 {
				e.writeSpace(buf)
			}
			e.writeAny(buf, key)
			buf.WriteByte(':')
			e.writeAny(buf, val.MapIndex(key))
		}
		buf.WriteByte(']')
	case reflect.Ptr:
		if val.IsNil() {
			buf.WriteString("nil")
			return
		}
		buf.WriteString("0x")
		e.writeUint(buf, uint64(val.Pointer()), 16)
	case reflect.Invalid:
		buf.WriteString("nil")
	}
}

func (e *encodeText) writeNewline(buf *buffer) {
	buf.WriteByte('\n')
}

func (e *encodeText) writeSpace(buf *buffer) {
	buf.WriteByte(' ')
}

func (e *encodeText) writeString(buf *buffer, s string) {
	buf.WriteByte('"')
	buf.WriteString(s)
	buf.WriteByte('"')
}

func (e *encodeText) writeBool(buf *buffer, b bool) {
	if b {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
}

func (e *encodeText) writeInt(buf *buffer, i int64) {
	*buf = strconv.AppendInt(*buf, i, 10)
}

func (e *encodeText) writeUint(buf *buffer, u uint64, base int) {
	*buf = strconv.AppendUint(*buf, u, base)
}

func (e *encodeText) writeFloat(buf *buffer, f float64) {
	*buf = strconv.AppendFloat(*buf, f, 'g', -1, 64)
}

func (e *encodeText) writeDuration(buf *buffer, d time.Duration) {
	*buf = strconv.AppendInt(*buf, int64(d), 10)
}

func (e *encodeText) writeTimeRFC3339(buf *buffer, t time.Time) {
	*buf = appendRFC3339Millis(*buf, t)
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
