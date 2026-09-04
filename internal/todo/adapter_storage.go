package todo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/serror"
)

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *storage {
	return &storage{db: db}
}

func (s *storage) HasUser(ctx context.Context, userID string) (bool, error) {
	var count int
	err := s.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users WHERE id = $1", userID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type todoRecord struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	UserID      string    `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
}

func (t todoRecord) Todo() Todo {
	return Todo{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Status:      t.Status,
		UserID:      t.UserID,
		CreatedAt:   t.CreatedAt,
	}
}

func (s *storage) GetTodos(ctx context.Context, userID string) ([]Todo, error) {
	var todos []todoRecord
	err := s.db.SelectContext(ctx, &todos, "SELECT * FROM todos WHERE user_id = ?", userID)
	if err != nil {
		return nil, serror.From(err)
	}
	var result []Todo
	for _, todo := range todos {
		result = append(result, todo.Todo())
	}
	return result, nil
}

func (s *storage) CreateTodo(ctx context.Context, todo Todo) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO todos (id, name, description, status, user_id, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		todo.ID, todo.Name, todo.Description, todo.Status, todo.UserID, todo.CreatedAt)
	if err != nil {
		return serror.From(err)
	}
	return nil
}

func (s *storage) FindTodo(ctx context.Context, userID, todoID string) (Todo, bool, error) {
	var todo todoRecord
	err := s.db.GetContext(ctx, &todo, "SELECT * FROM todos WHERE id = ? AND user_id = ?", todoID, userID)
	if err == sql.ErrNoRows {
		return Todo{}, false, nil
	}
	if err != nil {
		return Todo{}, false, serror.From(err)
	}
	return todo.Todo(), true, nil
}

func (s *storage) UpdateTodo(ctx context.Context, todo Todo) error {
	_, err := s.db.ExecContext(ctx, "UPDATE todos SET name = ?, description = ?, status = ? WHERE id = ?",
		todo.Name, todo.Description, todo.Status, todo.ID)
	if err != nil {
		return serror.From(err)
	}
	return nil
}

func (s *storage) DeleteTodo(ctx context.Context, todoID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM todos WHERE id = ?", todoID)
	if err != nil {
		return serror.From(err)
	}
	return nil
}
