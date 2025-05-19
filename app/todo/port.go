package todo

import "github.com/kongsakchai/gotemplate/app"

type Servicer interface {
	Todos(ctx app.Context) ([]Todo, error)
	Todo(ctx app.Context, id string) (Todo, error)
	Create(ctx app.Context, todo Todo) (Todo, error)
	Update(ctx app.Context, todo Todo) (Todo, error)
	Delete(ctx app.Context, id string) error
}
