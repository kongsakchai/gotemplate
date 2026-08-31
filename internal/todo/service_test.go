package todo

import (
	"context"
	"errors"
	"testing"

	"github.com/kongsakchai/gotemplate/pkg/serror"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodo(t *testing.T) {
	t.Run("should create todo when user exists", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		inputTodo := Todo{Name: "Test Todo", Status: "pending"}

		mock.EXPECT().HasUser(ctx, userID).Return(true, nil)
		expectedTodo := Todo{Name: "Test Todo", Status: "pending", UserID: userID}
		mock.EXPECT().CreateTodo(ctx, expectedTodo).Return(nil)

		err := svc.CreateTodo(ctx, userID, inputTodo)
		assert.NoError(t, err)
	})

	t.Run("should return ErrUserNotFound when user does not exist", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-2"
		todo := Todo{Name: "Test Todo", Status: "pending"}

		mock.EXPECT().HasUser(ctx, userID).Return(false, nil)

		err := svc.CreateTodo(ctx, userID, todo)

		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, serr.Code(), ErrUserNotFound.Code)
		assert.Equal(t, serr.Msg(), ErrUserNotFound.Msg)
	})

	t.Run("should propagate storage error from HasUser", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-3"
		todo := Todo{Name: "Test Todo", Status: "pending"}
		dbErr := errors.New("database error")

		mock.EXPECT().HasUser(ctx, userID).Return(false, dbErr)

		err := svc.CreateTodo(ctx, userID, todo)
		assert.Same(t, dbErr, err)
	})
}

func TestGetTodos(t *testing.T) {
	t.Run("should return todos for valid user", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		expectedTodos := []Todo{
			{ID: "1", Name: "Todo 1", UserID: "user-1", Status: "pending"},
			{ID: "2", Name: "Todo 2", UserID: "user-1", Status: "done"},
		}

		mock.EXPECT().GetTodos(ctx, userID).Return(expectedTodos, nil)

		todos, err := svc.GetTodos(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, todos, len(expectedTodos))
	})

	t.Run("should propagate storage error", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		dbErr := errors.New("database error")

		mock.EXPECT().GetTodos(ctx, userID).Return(nil, dbErr)

		_, err := svc.GetTodos(ctx, userID)
		assert.Same(t, dbErr, err)
	})
}

func TestUpdateTodo(t *testing.T) {
	t.Run("should update todo successfully", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		existingTodo := Todo{ID: "1", Name: "Old Name", Status: "pending", UserID: "user-1"}
		inputTodo := Todo{ID: "1", Name: "New Name", Status: "done"}

		mock.EXPECT().FindTodo(ctx, userID, existingTodo.ID).Return(existingTodo, true, nil)
		expectedTodo := Todo{ID: "1", Name: "New Name", Status: "done", UserID: userID}
		mock.EXPECT().UpdateTodo(ctx, expectedTodo).Return(nil)

		err := svc.UpdateTodo(ctx, userID, inputTodo)
		assert.NoError(t, err)
	})

	t.Run("should return ErrTodoNotFound when todo not found", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		updateTodo := Todo{ID: "1", Name: "New Name"}

		mock.EXPECT().FindTodo(ctx, userID, "1").Return(Todo{}, false, nil)

		err := svc.UpdateTodo(ctx, userID, updateTodo)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, serr.Code(), ErrTodoNotFound.Code)
		assert.Equal(t, serr.Msg(), ErrTodoNotFound.Msg)
	})

	t.Run("should propagate storage error from FindTodo", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		updateTodo := Todo{ID: "1", Name: "New Name"}
		dbErr := errors.New("database error")

		mock.EXPECT().FindTodo(ctx, userID, "1").Return(Todo{}, false, dbErr)

		err := svc.UpdateTodo(ctx, userID, updateTodo)
		assert.Same(t, dbErr, err)
	})
}

func TestDeleteTodo(t *testing.T) {
	t.Run("should delete todo successfully", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		todoID := "1"
		existingTodo := Todo{ID: todoID, Name: "Todo", Status: "pending", UserID: userID}

		mock.EXPECT().FindTodo(ctx, userID, todoID).Return(existingTodo, true, nil)
		mock.EXPECT().DeleteTodo(ctx, todoID).Return(nil)

		err := svc.DeleteTodo(ctx, userID, todoID)
		assert.NoError(t, err)
	})

	t.Run("should return ErrTodoNotFound when todo not found", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		todoID := "1"

		mock.EXPECT().FindTodo(ctx, userID, todoID).Return(Todo{}, false, nil)

		err := svc.DeleteTodo(ctx, userID, todoID)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, serr.Code(), ErrTodoNotFound.Code)
		assert.Equal(t, serr.Msg(), ErrTodoNotFound.Msg)
	})

	t.Run("should propagate storage error from FindTodo", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		todoID := "1"
		dbErr := errors.New("database error")

		mock.EXPECT().FindTodo(ctx, userID, todoID).Return(Todo{}, false, dbErr)

		err := svc.DeleteTodo(ctx, userID, todoID)
		assert.Same(t, dbErr, err)
	})

	t.Run("should propagate storage error from DeleteTodo", func(t *testing.T) {
		mock := newMockStorager(t)
		svc := NewService(mock)

		ctx := context.Background()
		userID := "user-1"
		todoID := "1"
		existingTodo := Todo{ID: todoID, Name: "Todo", UserID: userID}
		dbErr := errors.New("delete failed")

		mock.EXPECT().FindTodo(ctx, userID, todoID).Return(existingTodo, true, nil)
		mock.EXPECT().DeleteTodo(ctx, todoID).Return(dbErr)

		err := svc.DeleteTodo(ctx, userID, todoID)
		assert.Same(t, dbErr, err)
	})
}
