package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Writer interface {
	WriteHeader(code int)
	Write(b []byte)
}

type ginWriter struct {
	gin.ResponseWriter
	custom Writer
}

func (w *ginWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.custom != nil {
		w.custom.WriteHeader(code)
	}
}

func (w *ginWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if w.custom != nil {
		w.custom.Write(b)
	}
	return n, err
}

type echoResponseWriter struct {
	http.ResponseWriter
	custom Writer
}

func (w *echoResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.custom != nil {
		w.custom.WriteHeader(code)
	}
}

func (w *echoResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if w.custom != nil {
		w.custom.Write(b)
	}
	return n, err
}
