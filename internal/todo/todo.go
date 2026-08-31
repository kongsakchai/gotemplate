package todo

import (
	"context"
	"time"

	"github.com/kongsakchai/gotemplate/pkg/serror"
)

var (
	ErrUserNotFound = serror.NewCoded("3000", "user not found")
	ErrTodoNotFound = serror.NewCoded("3001", "todo not found")
)

type Todo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	UserID      string    `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
}

//mockery:generate: true
type Storager interface {
	HasUser(ctx context.Context, userID string) (bool, error)
	CreateTodo(ctx context.Context, todo Todo) error
	GetTodos(ctx context.Context, userID string) ([]Todo, error)
	FindTodo(ctx context.Context, userID string, id string) (Todo, bool, error)
	UpdateTodo(ctx context.Context, todo Todo) error
	DeleteTodo(ctx context.Context, id string) error
}

type Servicer interface {
	GetTodos(ctx context.Context, userID string) ([]Todo, error)
	CreateTodo(ctx context.Context, userID string, todo Todo) error
	UpdateTodo(ctx context.Context, userID string, todo Todo) error
	DeleteTodo(ctx context.Context, userID string, id string) error
}
