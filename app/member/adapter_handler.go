package member

import (
	"errors"
	"time"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
)

type handler struct {
	service Servicer
}

func NewHandler(service Servicer) *handler {
	return &handler{service: service}
}

func (h *handler) RegisterMemberHandler(app *app.EchoApp) {
	api := app.Group("/api/v1/members")
	api.GET("/", h.members)
	api.GET("/:username", h.member)
	api.POST("/", h.create)
	api.PUT("/:username", h.update)
	api.DELETE("/:username", h.remove)
}

func (h *handler) handlerError(err error) error {
	switch {
	case errors.Is(err, ErrorMaxAge) || errors.Is(err, ErrorMinAge):
		return app.BadRequest(app.InvalidAgeCode, app.InvalidAgeMsg, err)
	case errors.Is(err, ErrorDuplicate):
		return app.Conflict(app.UsernameUnavailableCode, app.UsernameUnavailableMsg, err)
	case errors.Is(err, ErrorMemberNotFound):
		return app.NotFound(app.MemberNotFoundCode, app.MemberNotFoundMsg, err)
	default:
		return app.InternalError(app.InternalErrorCode, app.InternalErrorMsg, err)
	}
}

func (h *handler) members(ctx *echo.Context) error {
	members, err := h.service.Members(ctx.Request().Context())
	if err != nil {
		return h.handlerError(err)
	}
	return app.Ok(ctx, members)
}

type usernameParam struct {
	Username string `param:"username" validate:"required"`
}

func (h *handler) member(ctx *echo.Context) error {
	req := usernameParam{}
	if err := ctx.Bind(&req); err != nil {
		return app.BadRequest(app.BadRequestCode, app.BadRequestMsg, err)
	}

	member, err := h.service.Member(ctx.Request().Context(), req.Username)
	if err != nil {
		return h.handlerError(err)
	}
	return app.Ok(ctx, member)
}

func (h *handler) remove(ctx *echo.Context) error {
	req := usernameParam{}
	if err := ctx.Bind(&req); err != nil {
		return app.BadRequest(app.BadRequestCode, app.BadRequestMsg, err)
	}

	if err := h.service.Remove(ctx.Request().Context(), req.Username); err != nil {
		return h.handlerError(err)
	}
	return app.Ok(ctx, nil)
}

type createBody struct {
	Username     string    `json:"username" validate:"required"`
	FirstName    string    `json:"firstName" validate:"required"`
	LastName     string    `json:"lastName" validate:"required"`
	Birthday     time.Time `json:"birthday" validate:"required"`
	RegisterDate time.Time `json:"registerDate"`
}

func (h *handler) create(ctx *echo.Context) error {
	req := createBody{}
	if err := app.Request(ctx, &req); err != nil {
		return err
	}
	if err := h.service.Create(ctx.Request().Context(), Member(req)); err != nil {
		return h.handlerError(err)
	}
	return app.Created(ctx, nil)
}

type updateBody struct {
	Username string `param:"username" validate:"required"`

	FirstName    string    `json:"firstName" validate:"required"`
	LastName     string    `json:"lastName" validate:"required"`
	Birthday     time.Time `json:"birthday" validate:"required"`
	RegisterDate time.Time `json:"registerDate"`
}

func (h *handler) update(ctx *echo.Context) error {
	req := updateBody{}
	if err := app.Request(ctx, &req); err != nil {
		return err
	}
	if err := h.service.Update(ctx.Request().Context(), req.Username, Member(req)); err != nil {
		return h.handlerError(err)
	}
	return app.Created(ctx, nil)
}
