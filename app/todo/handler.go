package todo

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/pkg/generate"
)

type handler struct {
	service Servicer
}

func NewHandler(service Servicer) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Todos(ctx app.Context) error {
	todos, err := h.service.Todos(ctx)
	if err != nil {
		return ctx.Error(app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, err))
	}

	return ctx.OK(todos)
}

func (h *handler) Todo(ctx app.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.Error(app.BadRequestError(app.MissingRequiredCode, app.MissingRequiredMsg))
	}

	todo, err := h.service.Todo(ctx, id)
	if err != nil {
		return ctx.Error(app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, err))
	}

	return ctx.OK(todo)
}

func (h *handler) Create(ctx app.Context) error {
	todo := Todo{}
	if err := ctx.Bind(&todo); err != nil {
		return ctx.Error(app.BadRequestError(app.InvalidRequestCode, app.InvalidRequestMsg, err))
	}

	todo.ID = generate.UUID()
	todo, err := h.service.Create(ctx, todo)
	if err != nil {
		return ctx.Error(app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, err))
	}

	return ctx.Created(todo)
}

func (h *handler) Update(ctx app.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.Error(app.BadRequestError(app.MissingRequiredCode, app.MissingRequiredMsg))
	}

	todo := Todo{}
	if err := ctx.Bind(&todo); err != nil {
		return ctx.Error(app.BadRequestError(app.InvalidRequestCode, app.InvalidRequestMsg, err))
	}

	todo.ID = id
	todo, err := h.service.Update(ctx, todo)
	if err != nil {
		return ctx.Error(app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, err))
	}

	return ctx.OK(todo)
}

func (h *handler) Delete(ctx app.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.Error(app.BadRequestError(app.MissingRequiredCode, app.MissingRequiredMsg))
	}

	if err := h.service.Delete(ctx, id); err != nil {
		return ctx.Error(app.InternalServerError(app.UnknownErrorCode, app.UnknowErrorMsg, err))
	}

	return ctx.OK(nil)
}
