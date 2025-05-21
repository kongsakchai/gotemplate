package todo

import "github.com/kongsakchai/gotemplate/app"

type service struct {
	todos []Todo
}

func NewService() *service {
	return &service{
		todos: []Todo{},
	}
}

func (s *service) Todos(ctx app.Context) ([]Todo, error) {
	return s.todos, nil
}

func (s *service) Todo(ctx app.Context, id string) (Todo, error) {
	for _, todo := range s.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return Todo{}, app.NotFoundError(app.NotFoundErrorCode, app.NotFoundMsg)
}

func (s *service) Create(ctx app.Context, todo Todo) (Todo, error) {
	s.todos = append(s.todos, todo)
	return todo, nil
}

func (s *service) Update(ctx app.Context, todo Todo) (Todo, error) {
	for i, t := range s.todos {
		if t.ID == todo.ID {
			s.todos[i] = todo
			return todo, nil
		}
	}
	return Todo{}, app.NotFoundError(app.NotFoundErrorCode, app.NotFoundMsg)
}

func (s *service) Delete(ctx app.Context, id string) error {
	for i, todo := range s.todos {
		if todo.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return nil
		}
	}
	return app.NotFoundError(app.NotFoundErrorCode, app.NotFoundMsg)
}
