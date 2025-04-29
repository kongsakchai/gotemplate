package ping

import (
	"time"

	"github.com/kongsakchai/gotemplate/app"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Ping(c app.Context) error {
	time.Sleep(10 * time.Second)
	return c.OK("Pong !!!")
}
