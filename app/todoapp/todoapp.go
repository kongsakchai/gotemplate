package todoapp

import (
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/todo"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
	"github.com/labstack/echo/v5"
)

//mockery:generate: true
type Servicer = todo.Servicer

//mockery:generate: true
type Verifier = jwttoken.Verifier

type todoApp struct {
	sv       Servicer
	verifier Verifier
}

type Deps struct {
	DB       *sqlx.DB
	Verifier Verifier
}

func NewApp(deps Deps) *todoApp {
	st := todo.NewStorage(deps.DB)
	sv := todo.NewService(st)
	return &todoApp{
		sv:       sv,
		verifier: deps.Verifier,
	}
}

func (a *todoApp) RegisterRoutes(echo *app.EchoApp) {
	authGroup := echo.Group("/api/v1/todos", app.AuthMiddleware(a.verifier), app.ErrorMiddleware(handleError))
	app.GET(authGroup, "get_todos", "", a.getTodos)
	app.POST(authGroup, "create_todo", "", a.createTodo)
	app.PUT(authGroup, "update_todo", "/:id", a.updateTodo)
	app.DELETE(authGroup, "delete_todo", "/:id", a.deleteTodo)
}

func (a *todoApp) getTodos(c *echo.Context) error {
	sub := c.Get("sub").(string)
	todos, err := a.sv.GetTodos(c.Request().Context(), sub)
	if err != nil {
		return err
	}
	return app.Ok(c, todos)
}

type createTodoRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (a *todoApp) createTodo(c *echo.Context) error {
	req, err := app.Request[createTodoRequest](c)
	if err != nil {
		return err
	}
	todo := todo.Todo{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	sub := c.Get("sub").(string)
	if err := a.sv.CreateTodo(c.Request().Context(), sub, todo); err != nil {
		return err
	}
	return app.Created(c, nil)
}

type updateTodoRequest struct {
	ID          string `param:"id" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (a *todoApp) updateTodo(c *echo.Context) error {
	req, err := app.Request[updateTodoRequest](c)
	if err != nil {
		return err
	}
	todo := todo.Todo{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	sub := c.Get("sub").(string)
	if err := a.sv.UpdateTodo(c.Request().Context(), sub, todo); err != nil {
		return err
	}
	return app.Ok(c, nil)
}

type deleteTodoRequest struct {
	ID string `param:"id" validate:"required"`
}

func (a *todoApp) deleteTodo(c *echo.Context) error {
	req, err := app.Request[deleteTodoRequest](c)
	if err != nil {
		return err
	}
	sub := c.Get("sub").(string)
	if err := a.sv.DeleteTodo(c.Request().Context(), sub, req.ID); err != nil {
		return err
	}
	return app.Ok(c, nil)
}
