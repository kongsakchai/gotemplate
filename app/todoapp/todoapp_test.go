package todoapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/todo"
	pkgconfig "github.com/kongsakchai/gotemplate/pkg/config"
	"github.com/kongsakchai/gotemplate/pkg/serror"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleError(t *testing.T) {
	t.Run("should return NotFound for ErrUserNotFound", func(t *testing.T) {
		appErr := todo.ErrUserNotFound.Err()
		result := handleError(appErr)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.Equal(t, todo.ErrUserNotFound.Code, result.Code)
		assert.Equal(t, todo.ErrUserNotFound.Msg, result.Message)
	})

	t.Run("should return NotFound for ErrTodoNotFound", func(t *testing.T) {
		appErr := todo.ErrTodoNotFound.Err()
		result := handleError(appErr)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.Equal(t, todo.ErrTodoNotFound.Code, result.Code)
		assert.Equal(t, todo.ErrTodoNotFound.Msg, result.Message)
	})

	t.Run("should return InternalError for unknown serror", func(t *testing.T) {
		customErr := serror.NewCoded("9997", "other error").Err()
		result := handleError(customErr)

		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
		assert.Equal(t, app.InternalErrorCode, result.Code)
		assert.Equal(t, app.InternalErrorMsg, result.Message)
	})

	t.Run("should return InternalError for non-serror", func(t *testing.T) {
		appErr := assert.AnError
		result := handleError(appErr)

		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
		assert.Equal(t, app.InternalErrorCode, result.Code)
		assert.Equal(t, app.InternalErrorMsg, result.Message)
	})

	t.Run("should return InternalError for nil error", func(t *testing.T) {
		result := handleError(nil)

		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
	})
}

func TestNewApp(t *testing.T) {
	t.Run("should create todoApp with dependencies", func(t *testing.T) {
		a := NewApp(Deps{})
		assert.NotNil(t, a)
	})
}

func TestTodoAppRouteRegistration(t *testing.T) {
	t.Run("should register routes without panic", func(t *testing.T) {
		cfg := pkgconfig.Config{}
		echoApp := app.NewEchoApp(cfg)
		a := NewApp(Deps{})

		assert.NotPanics(t, func() {
			a.RegisterRoutes(echoApp)
		})
	})
}

func TestTodoAppHandlers(t *testing.T) {
	// newContext creates an echotest context with an optional JSON body and
	// path values, and injects the authenticated user id (normally set by
	// app.AuthMiddleware) into the context.
	newContext := func(body string, pathValues echo.PathValues) (*echo.Context, *httptest.ResponseRecorder) {
		ctx, rec := echotest.ContextConfig{
			JSONBody:   []byte(body),
			PathValues: pathValues,
		}.ToContextRecorder(t)
		ctx.Set("sub", "user-1")
		return ctx, rec
	}

	// === getTodos handler ===

	t.Run("getTodos - success returns 200 OK with todos list", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().GetTodos(mock.Anything, "user-1").Return([]todo.Todo{
			{ID: "1", Name: "Buy groceries", Status: "pending"},
			{ID: "2", Name: "Call doctor", Status: "urgent"},
		}, nil)

		a := &todoApp{sv: sv}
		ctx, rec := newContext("", nil)
		err := a.getTodos(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp app.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
	})

	t.Run("getTodos - service error returns error", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().GetTodos(mock.Anything, "user-1").Return(nil, assert.AnError)

		a := &todoApp{sv: sv}
		ctx, _ := newContext("", nil)
		err := a.getTodos(ctx)

		assert.ErrorIs(t, err, assert.AnError)
	})

	// === createTodo handler ===

	t.Run("createTodo - success returns 201 Created", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().CreateTodo(mock.Anything, "user-1", mock.Anything).Return(nil)

		a := &todoApp{sv: sv}
		ctx, rec := newContext(`{"name":"New task","description":"Do something","status":"pending"}`, nil)
		err := a.createTodo(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp app.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
	})

	t.Run("createTodo - invalid JSON returns bind error", func(t *testing.T) {
		a := &todoApp{}
		ctx, _ := newContext(`{invalid json}`, nil)
		err := a.createTodo(ctx)

		assert.Error(t, err)
	})

	t.Run("createTodo - service error returns serror", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().CreateTodo(mock.Anything, "user-1", mock.Anything).Return(todo.ErrUserNotFound.Err())

		a := &todoApp{sv: sv}
		ctx, _ := newContext(`{"name":"New task"}`, nil)
		err := a.createTodo(ctx)

		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, todo.ErrUserNotFound.Code, serr.Code())
		assert.Equal(t, todo.ErrUserNotFound.Msg, serr.Msg())
	})

	// === updateTodo handler ===

	t.Run("updateTodo - success returns 200 OK", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().UpdateTodo(mock.Anything, "user-1", mock.Anything).Return(nil)

		a := &todoApp{sv: sv}
		ctx, rec := newContext(`{"name":"Updated task","description":"Updated desc","status":"done"}`,
			echo.PathValues{{Name: "id", Value: "todo-update-1"}})
		err := a.updateTodo(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp app.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
	})

	t.Run("updateTodo - invalid JSON returns bind error", func(t *testing.T) {
		a := &todoApp{}
		ctx, _ := newContext(`bad json`, echo.PathValues{{Name: "id", Value: "todo-update-1"}})
		err := a.updateTodo(ctx)

		assert.Error(t, err)
	})

	t.Run("updateTodo - service error returns serror", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().UpdateTodo(mock.Anything, "user-1", mock.Anything).Return(todo.ErrTodoNotFound.Err())

		a := &todoApp{sv: sv}
		ctx, _ := newContext(`{"name":"Missing"}`, echo.PathValues{{Name: "id", Value: "missing-id"}})
		err := a.updateTodo(ctx)

		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, todo.ErrTodoNotFound.Code, serr.Code())
		assert.Equal(t, todo.ErrTodoNotFound.Msg, serr.Msg())
	})

	// === deleteTodo handler ===

	t.Run("deleteTodo - success returns 200 OK", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().DeleteTodo(mock.Anything, "user-1", "todo-delete-1").Return(nil)

		a := &todoApp{sv: sv}
		ctx, rec := newContext("", echo.PathValues{{Name: "id", Value: "todo-delete-1"}})
		err := a.deleteTodo(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp app.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
	})

	t.Run("deleteTodo - invalid JSON returns bind error", func(t *testing.T) {
		a := &todoApp{}
		ctx, _ := newContext(`{invalid json}`, echo.PathValues{{Name: "id", Value: "todo-delete-1"}})
		err := a.deleteTodo(ctx)

		assert.Error(t, err)
	})

	t.Run("deleteTodo - service error returns serror", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().DeleteTodo(mock.Anything, "user-1", "missing-todo").Return(todo.ErrTodoNotFound.Err())

		a := &todoApp{sv: sv}
		ctx, _ := newContext("", echo.PathValues{{Name: "id", Value: "missing-todo"}})
		err := a.deleteTodo(ctx)

		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, todo.ErrTodoNotFound.Code, serr.Code())
		assert.Equal(t, todo.ErrTodoNotFound.Msg, serr.Msg())
	})
}
