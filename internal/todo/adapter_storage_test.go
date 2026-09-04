package todo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/serror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	t.Run("should create storage with db", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		assert.NotNil(t, st)
		assert.NotNil(t, st.db)
	})
}

func TestStorage_HasUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should return true when user exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id").
			WithArgs("user-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		exists, err := st.HasUser(ctx, "user-1")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when user does not exist", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id").
			WithArgs("user-2").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		exists, err := st.HasUser(ctx, "user-2")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should propagate database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id").
			WithArgs("user-3").
			WillReturnError(assert.AnError)

		exists, err := st.HasUser(ctx, "user-3")
		assert.Error(t, err)
		assert.False(t, exists)
	})
}

func TestStorage_GetTodos(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("should return todos for valid user", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		rows := sqlmock.NewRows([]string{"id", "name", "description", "status", "user_id", "created_at"}).
			AddRow("1", "Todo 1", "Desc 1", "pending", "user-1", now).
			AddRow("2", "Todo 2", "Desc 2", "done", "user-1", now)

		mock.ExpectQuery("SELECT \\* FROM todos WHERE user_id").
			WithArgs("user-1").
			WillReturnRows(rows)

		todos, err := st.GetTodos(ctx, "user-1")
		assert.NoError(t, err)
		assert.Len(t, todos, 2)
		assert.Equal(t, "1", todos[0].ID)
		assert.Equal(t, "Todo 1", todos[0].Name)
		assert.Equal(t, "pending", todos[0].Status)
	})

	t.Run("should return empty slice when no todos found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT \\* FROM todos WHERE user_id").
			WithArgs("user-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "status", "user_id", "created_at"}))

		todos, err := st.GetTodos(ctx, "user-1")
		assert.NoError(t, err)
		assert.Empty(t, todos)
	})

	t.Run("should wrap and return database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT \\* FROM todos WHERE user_id").
			WithArgs("user-1").
			WillReturnError(assert.AnError)

		todos, err := st.GetTodos(ctx, "user-1")
		assert.Nil(t, todos)
		assert.NotNil(t, err)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		_ = serr // serror is wrapped
	})
}

func TestStorage_CreateTodo(t *testing.T) {
	ctx := context.Background()

	t.Run("should create todo successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("INSERT INTO todos").
			WithArgs("todo-1", "Test Todo", "Desc", "pending", "user-1", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = st.CreateTodo(ctx, Todo{
			ID:          "todo-1",
			Name:        "Test Todo",
			Description: "Desc",
			Status:      "pending",
			UserID:      "user-1",
		})
		assert.NoError(t, err)
	})

	t.Run("should wrap and return database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("INSERT INTO todos").
			WithArgs("todo-1", "Test Todo", "Desc", "pending", "user-1", sqlmock.AnyArg()).
			WillReturnError(assert.AnError)

		err = st.CreateTodo(ctx, Todo{
			ID:     "todo-1",
			Name:   "Test Todo",
			Status: "pending",
			UserID: "user-1",
		})
		assert.NotNil(t, err)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		_ = serr
	})
}

func TestStorage_FindTodo(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("should return todo when found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		rows := sqlmock.NewRows([]string{"id", "name", "description", "status", "user_id", "created_at"}).
			AddRow("1", "Todo 1", "Desc 1", "pending", "user-1", now)

		mock.ExpectQuery("SELECT \\* FROM todos WHERE id = .+ AND user_id = .+").
			WithArgs("1", "user-1").
			WillReturnRows(rows)

		todo, found, err := st.FindTodo(ctx, "user-1", "1")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "1", todo.ID)
		assert.Equal(t, "Todo 1", todo.Name)
	})

	t.Run("should return false when todo not found (ErrNoRows)", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT \\* FROM todos WHERE id = .+ AND user_id = .+").
			WithArgs("1", "user-1").
			WillReturnError(sql.ErrNoRows)

		todo, found, err := st.FindTodo(ctx, "user-1", "1")
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Equal(t, Todo{}, todo)
	})

	t.Run("should wrap and return database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT \\* FROM todos WHERE id = .+ AND user_id = .+").
			WithArgs("1", "user-1").
			WillReturnError(assert.AnError)

		todo, found, err := st.FindTodo(ctx, "user-1", "1")
		assert.Error(t, err)
		assert.False(t, found)
		assert.Equal(t, Todo{}, todo)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		_ = serr
	})
}

func TestStorage_UpdateTodo(t *testing.T) {
	ctx := context.Background()

	t.Run("should update todo successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("UPDATE todos SET name = .+, description = .+, status = .+ WHERE id = .+").
			WithArgs("New Name", "New Desc", "done", "1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = st.UpdateTodo(ctx, Todo{
			ID:          "1",
			Name:        "New Name",
			Description: "New Desc",
			Status:      "done",
		})
		assert.NoError(t, err)
	})

	t.Run("should wrap and return database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("UPDATE todos SET name = .+, description = .+, status = .+ WHERE id = .+").
			WithArgs("New Name", "New Desc", "done", "1").
			WillReturnError(assert.AnError)

		err = st.UpdateTodo(ctx, Todo{
			ID:          "1",
			Name:        "New Name",
			Description: "New Desc",
			Status:      "done",
		})
		assert.NotNil(t, err)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		_ = serr
	})
}

func TestStorage_DeleteTodo(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete todo successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("DELETE FROM todos WHERE id = .+").
			WithArgs("1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = st.DeleteTodo(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("should wrap and return database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("DELETE FROM todos WHERE id = .+").
			WithArgs("1").
			WillReturnError(assert.AnError)

		err = st.DeleteTodo(ctx, "1")
		assert.NotNil(t, err)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		_ = serr
	})
}

func TestTodoRecord_Todo(t *testing.T) {
	t.Run("should convert todoRecord to Todo", func(t *testing.T) {
		now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		record := todoRecord{
			ID:          "1",
			Name:        "Test Todo",
			Description: "Description",
			Status:      "pending",
			UserID:      "user-1",
			CreatedAt:   now,
		}

		todo := record.Todo()

		assert.Equal(t, "1", todo.ID)
		assert.Equal(t, "Test Todo", todo.Name)
		assert.Equal(t, "Description", todo.Description)
		assert.Equal(t, "pending", todo.Status)
		assert.Equal(t, "user-1", todo.UserID)
		assert.Equal(t, now, todo.CreatedAt)
	})
}
