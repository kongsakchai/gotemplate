package member

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/domain"
	"github.com/labstack/echo/v5"
)

type Deps struct {
	DB    *sqlx.DB
	Clock Clock
}

type memberHandler struct {
	service    Servicer
	badRequest map[string]struct{}
	conflict   map[string]struct{}
	notFound   map[string]struct{}
}

func NewMemberHandler(deps Deps) *memberHandler {
	badRequest := map[string]struct{}{InvalidAgeCode: {}}
	conflict := map[string]struct{}{UsernameUnavailableCode: {}}
	notFound := map[string]struct{}{MemberNotFoundCode: {}}

	st := NewStorage(deps.DB)
	sv := NewService(ServiceDeps{
		Clock:   deps.Clock,
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
	if d, ok := err.(domain.Error); ok {
		if _, ok := h.badRequest[d.Code]; ok {
			return app.BadRequest(d.Code, d.Message, err)
		}
		if _, ok := h.conflict[d.Code]; ok {
			return app.Conflict(d.Code, d.Message, err)
		}
		if _, ok := h.notFound[d.Code]; ok {
			return app.NotFound(d.Code, d.Message, err)
		}
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
	req, err := app.Request[memberReq](ctx)
	if err != nil {
		return err
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
	req, err := app.Request[removeReq](ctx)
	if err != nil {
		return err
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
	req, err := app.Request[createReq](ctx)
	if err != nil {
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
	req, err := app.Request[updateReq](ctx)
	if err != nil {
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
