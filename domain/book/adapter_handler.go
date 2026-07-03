package book

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
)

type Option struct {
	DB *sqlx.DB
}

type memberHandler struct {
	service  Servicer
	notFound app.ErrorMaps
}

func NewMemberHandler(option Option) *memberHandler {
	notFound := app.ErrorMaps{
		{Cause: ErrBookNotFound, Code: app.MemberNotFoundCode, Message: app.MemberNotFoundMsg},
	}

	st := NewStorage(option.DB)
	sv := NewService(ServiceOption{
		Clock:   option.Clock,
		Storage: st,
	})

	return &memberHandler{
		service:    sv,
		badRequest: badRequest,
		conflict:   conflict,
		notFound:   notFound,
	}
}

func (h *memberHandler) RegisterMemberHandler(e *echo.Echo) {
	apiPath := e.Group("/api/v1/members")
	apiPath.GET("", h.members, app.TagMiddleware("get_members"))
	apiPath.GET("/:username", h.member, app.TagMiddleware("get_member"))
	apiPath.POST("", h.create, app.TagMiddleware("create_member"))
	apiPath.PUT("/:username", h.update, app.TagMiddleware("update_member"))
	apiPath.DELETE("/:username", h.remove, app.TagMiddleware("remove_members"))
}

func (h *memberHandler) handlerError(err error) error {
	if mapping, ok := h.badRequest.Resolve(err); ok {
		return app.BadRequest(mapping.Code, mapping.Message, err)
	}
	if mapping, ok := h.conflict.Resolve(err); ok {
		return app.Conflict(mapping.Code, mapping.Message, err)
	}
	if mapping, ok := h.notFound.Resolve(err); ok {
		return app.NotFound(mapping.Code, mapping.Message, err)
	}

	return app.InternalError(app.InternalErrorCode, app.InternalErrorMsg, err)
}

func (h *memberHandler) members(ctx *echo.Context) error {
	members, err := h.service.GetMembers(ctx.Request().Context())
	if err != nil {
		return h.handlerError(err)
	}
	return app.Ok(ctx, members)
}

type memberReq struct {
	Username string `param:"username" validate:"required"`
}

func (h *memberHandler) member(ctx *echo.Context) error {
	req := memberReq{}
	if err := ctx.Bind(&req); err != nil {
		return app.BadRequest(app.BadRequestCode, app.BadRequestMsg, err)
	}

	member, err := h.service.GetMember(ctx.Request().Context(), req.Username)
	if err != nil {
		return h.handlerError(err)
	}
	return app.Ok(ctx, member)
}

type removeReq struct {
	Username string `param:"username" validate:"required"`
}

func (h *memberHandler) remove(ctx *echo.Context) error {
	req := removeReq{}
	if err := ctx.Bind(&req); err != nil {
		return app.BadRequest(app.BadRequestCode, app.BadRequestMsg, err)
	}

	if err := h.service.RemoveMember(ctx.Request().Context(), req.Username); err != nil {
		return h.handlerError(err)
	}
	return app.Ok(ctx, nil)
}

type createReq struct {
	Username  string    `json:"username" validate:"required"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Birthday  time.Time `json:"birthday" validate:"required"`
}

func (h *memberHandler) create(ctx *echo.Context) error {
	req := createReq{}
	if err := app.Request(ctx, &req); err != nil {
		return err
	}
	if err := h.service.CreateMember(ctx.Request().Context(), CreateMemberPayload{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  req.Birthday,
	}); err != nil {
		return h.handlerError(err)
	}
	return app.Created(ctx, nil)
}

type updateReq struct {
	Username string `param:"username" validate:"required"`

	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Birthday  time.Time `json:"birthday" validate:"required"`
}

func (h *memberHandler) update(ctx *echo.Context) error {
	req := updateReq{}
	if err := app.Request(ctx, &req); err != nil {
		return err
	}
	if err := h.service.UpdateMember(ctx.Request().Context(), UpdateMemberPayload{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  req.Birthday,
	}); err != nil {
		return h.handlerError(err)
	}
	return app.Created(ctx, nil)
}
