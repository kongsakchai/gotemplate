package todo

import (
	"errors"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/pkg/generate"
	"github.com/kongsakchai/gotemplate/test/appmock"
	"github.com/stretchr/testify/mock"
)

func TestHandlerTodos(t *testing.T) {
	t.Run("should return todos when service returns todos", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		todos := []Todo{
			{ID: "1", Title: "Todo 1", Completed: false},
			{ID: "2", Title: "Todo 2", Completed: true},
		}

		service.On("Todos", ctx).Return(todos, nil)
		ctx.On("OK", todos).Return(nil)

		// Act
		handler.Todos(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return error when service returns error", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedErr := errors.New("service error")

		service.On("Todos", ctx).Return([]Todo{}, expectedErr)
		ctx.On("Error", app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, expectedErr)).Return(nil)

		// Act
		handler.Todos(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})
}

func TestHandlerTodo(t *testing.T) {
	t.Run("should return todo when service returns todo", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		id := "1"
		todo := Todo{ID: id, Title: "Todo 1", Completed: false}

		ctx.On("Param", "id").Return("1")
		service.On("Todo", ctx, id).Return(todo, nil)
		ctx.On("OK", todo).Return(nil)

		// Act
		handler.Todo(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return error when service returns error", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedErr := errors.New("service error")

		ctx.On("Param", "id").Return("1")
		service.On("Todo", ctx, "1").Return(Todo{}, expectedErr)
		ctx.On("Error", app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, expectedErr)).Return(nil)

		// Act
		handler.Todo(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return bad request when id is missing", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		ctx.On("Param", "id").Return("")
		ctx.On("Error", app.BadRequestError(app.MissingRequiredCode, app.MissingRequiredMsg)).Return(nil)

		// Act
		handler.Todo(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})
}

func TestHandlerCreate(t *testing.T) {
	generate.SetFixedUUID("123456")
	t.Run("should create todo when service returns todo", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedTodo := Todo{ID: "123456", Title: "Todo 1", Completed: false}

		ctx.On("Bind", mock.AnythingOfType("*todo.Todo")).Return(nil).Run(appmock.ReturnArgs(expectedTodo))
		service.On("Create", ctx, expectedTodo).Return(expectedTodo, nil)
		ctx.On("Created", expectedTodo).Return(nil)

		// Act
		handler.Create(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return error when service returns error", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedTodo := Todo{ID: "123456", Title: "Todo 1", Completed: false}
		expectedErr := errors.New("service error")

		ctx.On("Bind", mock.AnythingOfType("*todo.Todo")).Return(nil).Run(appmock.ReturnArgs(expectedTodo))
		service.On("Create", ctx, expectedTodo).Return(Todo{}, expectedErr)
		ctx.On("Error", app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, expectedErr)).Return(nil)

		// Act
		handler.Create(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return bad request when bind fails", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedErr := errors.New("bind error")

		ctx.On("Bind", mock.AnythingOfType("*todo.Todo")).Return(expectedErr)
		ctx.On("Error", app.BadRequestError(app.InvalidRequestCode, app.InvalidRequestMsg, expectedErr)).Return(nil)

		// Act
		handler.Create(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})
}

func TestHandlerUpdate(t *testing.T) {
	t.Run("should update todo when service returns todo", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		id := "1"
		expectedTodo := Todo{ID: id, Title: "Todo 1", Completed: false}

		ctx.On("Param", "id").Return("1")
		ctx.On("Bind", mock.AnythingOfType("*todo.Todo")).Return(nil).Run(appmock.ReturnArgs(expectedTodo))
		service.On("Update", ctx, expectedTodo).Return(expectedTodo, nil)
		ctx.On("OK", expectedTodo).Return(nil)

		// Act
		handler.Update(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return error when service returns error", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedTodo := Todo{ID: "1", Title: "Todo 1", Completed: false}
		expectedErr := errors.New("service error")

		ctx.On("Param", "id").Return("1")
		ctx.On("Bind", mock.AnythingOfType("*todo.Todo")).Return(nil).Run(appmock.ReturnArgs(expectedTodo))
		service.On("Update", ctx, expectedTodo).Return(Todo{}, expectedErr)
		ctx.On("Error", app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, expectedErr)).Return(nil)

		// Act
		handler.Update(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return bad request when bind fails", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedErr := errors.New("bind error")

		ctx.On("Param", "id").Return("1")
		ctx.On("Bind", mock.AnythingOfType("*todo.Todo")).Return(expectedErr)
		ctx.On("Error", app.BadRequestError(app.InvalidRequestCode, app.InvalidRequestMsg, expectedErr)).Return(nil)

		// Act
		handler.Update(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return bad request when id is missing", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		ctx.On("Param", "id").Return("")
		ctx.On("Error", app.BadRequestError(app.MissingRequiredCode, app.MissingRequiredMsg)).Return(nil)

		// Act
		handler.Update(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

}

func TestHandlerDelete(t *testing.T) {
	t.Run("should delete todo when service returns no error", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		ctx.On("Param", "id").Return("1")
		service.On("Delete", ctx, "1").Return(nil)
		ctx.On("OK", nil).Return(nil)

		// Act
		handler.Delete(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return error when service returns error", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		expectedErr := errors.New("service error")

		ctx.On("Param", "id").Return("1")
		service.On("Delete", ctx, "1").Return(expectedErr)
		ctx.On("Error", app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, expectedErr)).Return(nil)

		// Act
		handler.Delete(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("should return bad request when id is missing", func(t *testing.T) {
		// Arrange
		service := new(mockService)
		handler := NewHandler(service)
		ctx := new(appmock.Context)

		ctx.On("Param", "id").Return("")
		ctx.On("Error", app.BadRequestError(app.MissingRequiredCode, app.MissingRequiredMsg)).Return(nil)

		// Act
		handler.Delete(ctx)

		// Assert
		service.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})
}
