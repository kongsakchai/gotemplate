package httpclient

import (
	"context"
	"net/http"
)

type OptionFunc func(*http.Request, context.Context) context.Context

func TraceOption(key string) OptionFunc {
	return func(r *http.Request, ctx context.Context) context.Context {
		trace, _ := ctx.Value(key).(string)
		r.Header.Set(key, trace)

		return ctx
	}
}
