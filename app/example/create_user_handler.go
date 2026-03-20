package example

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Age       int    `json:"age" validate:"required,gte=0,lte=130"`
}

func (h *handler) CreateUser(ctx *echo.Context) error {
	var req CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return app.BadRequest("4001", "invalid request body", err)
	}
	if err := ctx.Validate(&req); err != nil {
		return app.BadRequest("4002", "validation error", err, err)
	}

	err := h.createUser(req)
	if !err.IsEmpty() {
		return err
	}
	return app.Ok(ctx, nil)
}

func (h *handler) createUser(req CreateUserRequest) app.Error {
	user, err := h.storage.UserByName(req.FirstName)
	if err != nil {
		return app.InternalError("5001", "failed to get user by name", err)
	}

	if user.FirstName != "" {
		return app.Conflict("4003", "user already exists", nil)
	}

	err = h.storage.CreateUser(User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
	})
	if err != nil {
		return app.InternalError("5002", "failed to create user", err)
	}
	return app.Error{}
}
