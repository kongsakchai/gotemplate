package todo

import (
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/test/appmock"
	"github.com/stretchr/testify/assert"
)

func TestServiceTodos(t *testing.T) {
	t.Run("should return todos when repository returns todos", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		todos := []Todo{
			{ID: "1", Title: "Todo 1", Completed: false},
			{ID: "2", Title: "Todo 2", Completed: true},
		}
		service.todos = todos

		// Act
		result, err := service.Todos(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, todos, result)
	})
}

func TestServiceTodo(t *testing.T) {
	t.Run("should return todo when repository returns todo", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		todo := Todo{ID: "1", Title: "Todo 1", Completed: false}
		service.todos = []Todo{todo}

		// Act
		result, err := service.Todo(ctx, "1")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, todo, result)
	})

	t.Run("should return error when todo not found", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		service.todos = []Todo{
			{ID: "1", Title: "Todo 1", Completed: false},
			{ID: "2", Title: "Todo 2", Completed: true},
		}
		expectedError := app.NotFoundError(app.NotFoundErrorCode, app.NotFoundMsg)

		// Act
		result, err := service.Todo(ctx, "3")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, Todo{}, result)
	})
}

func TestServiceCreate(t *testing.T) {
	t.Run("should create todo when valid todo is provided", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		todo := Todo{ID: "1", Title: "Todo 1", Completed: false}

		// Act
		result, err := service.Create(ctx, todo)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, todo, result)
		assert.Len(t, service.todos, 1)
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("should update todo when valid todo is provided", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		todo := Todo{ID: "1", Title: "Todo 1", Completed: false}
		service.todos = []Todo{todo}

		updatedTodo := Todo{ID: "1", Title: "Updated Todo", Completed: true}

		// Act
		result, err := service.Update(ctx, updatedTodo)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updatedTodo, result)
		assert.Len(t, service.todos, 1)
	})

	t.Run("should return error when todo not found", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		todo := Todo{ID: "1", Title: "Todo 1", Completed: false}
		service.todos = []Todo{todo}

		updatedTodo := Todo{ID: "2", Title: "Updated Todo", Completed: true}
		expectedError := app.NotFoundError(app.NotFoundErrorCode, app.NotFoundMsg)

		// Act
		result, err := service.Update(ctx, updatedTodo)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, Todo{}, result)
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("should delete todo when valid id is provided", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		todo := Todo{ID: "1", Title: "Todo 1", Completed: false}
		service.todos = []Todo{todo}

		// Act
		err := service.Delete(ctx, "1")

		// Assert
		assert.NoError(t, err)
		assert.Len(t, service.todos, 0)
	})

	t.Run("should return error when todo not found", func(t *testing.T) {
		// Arrange
		service := NewService()
		ctx := new(appmock.Context)

		todo := Todo{ID: "1", Title: "Todo 1", Completed: false}
		service.todos = []Todo{todo}

		expectedError := app.NotFoundError(app.NotFoundErrorCode, app.NotFoundMsg)

		// Act
		err := service.Delete(ctx, "2")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}
