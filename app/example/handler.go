package example

import "github.com/kongsakchai/gotemplate/app"

//mockery:generate: true
type Storager interface {
	Users() ([]User, app.Error)
	UserByName(name string) (User, app.Error)
}

type handler struct {
	storage Storager
}

func NewHandler(storage Storager) *handler {
	return &handler{storage: storage}
}
