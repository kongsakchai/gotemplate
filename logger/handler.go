package logger

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"
)

type replaceAttrFunc func(groups []string, a slog.Attr) slog.Attr

type HandlerOptions struct {
	Level        slog.Level
	ReplaceAttr  replaceAttrFunc
	TimeFormat   string
	DisableColor bool
	DisableEmoji bool
	DisableTime  bool
	DisableLevel bool
}

type textHandler struct {
	opts       HandlerOptions
	attrPrefix []byte
	groups     []string

	mu sync.Mutex
	w  io.Writer
	en *encodeText
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
		groups: make([]string, 0, 5),
		en:     newEncodeText(*opts),
	}
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	buf := newBuffer()
	for _, a := range attrs {
		h.en.writeAttr(buf, h.groups, a)
	}
	return &textHandler{
		opts:       h.opts,
		w:          h.w,
		groups:     h.groups,
		en:         h.en,
		attrPrefix: *buf,
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
		en:         h.en,
		groups:     gs,
	}
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := newBuffer()
	defer buf.Free()

	// Write Emoji Level
	h.en.writeEmojiLevel(buf, r.Level)

	// Write Time
	h.en.writeTime(buf, r.Time)

	// Write Level
	h.en.writeLevel(buf, r.Level)

	// Write Message
	h.en.writeMessage(buf, r.Message)

	// Wrote attrPrefix
	if prefix := h.attrPrefix; len(prefix) > 0 {
		buf.Write(h.attrPrefix)
	}

	// Write Attributes
	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			if rep := h.opts.ReplaceAttr; rep != nil {
				a = rep(h.groups, a)
			}

			h.en.writeAttr(buf, h.groups, a)
			return true
		})
	}

	h.en.writeNewline(buf)

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*buf)
	return err
}
