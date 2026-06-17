package member

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kongsakchai/gotemplate/template/app"
	"github.com/kongsakchai/gotemplate/template/pkg/validator"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	t.Run("should create handler", func(t *testing.T) {
		svc := newMockServicer(t)
		h := NewHandler(svc)
		assert.NotNil(t, h)
	})

	t.Run("should register routes", func(t *testing.T) {
		svc := newMockServicer(t)
		h := NewHandler(svc)

		e := echo.New()
		app := &app.EchoApp{Echo: e}
		h.RegisterMemberHandler(app)

		routes := app.Router().Routes()
		assert.Len(t, routes, 5)
	})
}

func TestHandlerError(t *testing.T) {
	svc := newMockServicer(t)
	h := NewHandler(svc)

	t.Run("should return bad request for min age error", func(t *testing.T) {
		err := h.handlerError(ErrorMinAge)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, app.InvalidAgeCode, appErr.Code)
	})

	t.Run("should return bad request for max age error", func(t *testing.T) {
		err := h.handlerError(ErrorMaxAge)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
	})

	t.Run("should return conflict for duplicate error", func(t *testing.T) {
		err := h.handlerError(ErrorDuplicate)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusConflict, appErr.HTTPCode)
		assert.Equal(t, app.UsernameUnavailableCode, appErr.Code)
	})

	t.Run("should return not found for member not found error", func(t *testing.T) {
		err := h.handlerError(ErrorMemberNotFound)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, app.MemberNotFoundCode, appErr.Code)
	})

	t.Run("should return internal error for unknown error", func(t *testing.T) {
		err := h.handlerError(errors.New("unknown"))
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.HTTPCode)
		assert.Equal(t, app.InternalErrorCode, appErr.Code)
	})
}

func TestHandlerMembers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()

		svc := newMockServicer(t)
		svc.On("Members", contextBackground()).Return([]Member{member}, nil)

		h := NewHandler(svc)
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/api/v1/members", nil),
		}.ToContextRecorder(t)

		err := h.members(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := newMockServicer(t)
		svc.On("Members", contextBackground()).Return(nil, errors.New("service err"))

		h := NewHandler(svc)
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/api/v1/members", nil),
		}.ToContextRecorder(t)

		err := h.members(ctx)
		assert.Error(t, err)
	})
}

func TestHandlerMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()

		svc := newMockServicer(t)
		svc.On("Member", contextBackground(), "john").Return(member, nil)

		h := NewHandler(svc)
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/api/v1/members/john", nil),
			PathValues: echo.PathValues{
				{Name: "username", Value: "john"},
			},
		}.ToContextRecorder(t)

		err := h.member(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("bind error - invalid body with json content type", func(t *testing.T) {
		svc := newMockServicer(t)
		h := NewHandler(svc)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/members/john", strings.NewReader("{invalid}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx, _ := echotest.ContextConfig{
			Request: req,
		}.ToContextRecorder(t)

		err := h.member(ctx)
		assert.Error(t, err)
	})

	t.Run("service error", func(t *testing.T) {
		svc := newMockServicer(t)
		svc.On("Member", contextBackground(), "john").Return(Member{}, ErrorMemberNotFound)

		h := NewHandler(svc)
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/api/v1/members/john", nil),
			PathValues: echo.PathValues{
				{Name: "username", Value: "john"},
			},
		}.ToContextRecorder(t)

		err := h.member(ctx)
		assert.Error(t, err)
	})
}

func TestHandlerCreate(t *testing.T) {
	v := validator.NewReqValidator()

	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()

		svc := newMockServicer(t)
		svc.On("Create", contextBackground(), member).Return(nil)

		h := NewHandler(svc)
		body := `{"username":"john","firstName":"John","lastName":"Doe","birthday":"2000-01-01T00:00:00Z"}`
		ctx, rec := echotest.ContextConfig{
			Request:  httptest.NewRequest(http.MethodPost, "/api/v1/members", strings.NewReader(body)),
			JSONBody: []byte(body),
		}.ToContextRecorder(t)
		ctx.Echo().Validator = v

		err := h.create(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("bind error - invalid json", func(t *testing.T) {
		svc := newMockServicer(t)
		h := NewHandler(svc)
		ctx, _ := echotest.ContextConfig{
			Request:  httptest.NewRequest(http.MethodPost, "/api/v1/members", strings.NewReader("{invalid}")),
			JSONBody: []byte("{invalid}"),
		}.ToContextRecorder(t)
		ctx.Echo().Validator = v

		err := h.create(ctx)
		assert.Error(t, err)
	})

	t.Run("service error", func(t *testing.T) {
		member, _ := newFixture()

		svc := newMockServicer(t)
		svc.On("Create", contextBackground(), member).Return(ErrorDuplicate)

		h := NewHandler(svc)
		body := `{"username":"john","firstName":"John","lastName":"Doe","birthday":"2000-01-01T00:00:00Z"}`
		ctx, _ := echotest.ContextConfig{
			Request:  httptest.NewRequest(http.MethodPost, "/api/v1/members", strings.NewReader(body)),
			JSONBody: []byte(body),
		}.ToContextRecorder(t)
		ctx.Echo().Validator = v

		err := h.create(ctx)
		assert.Error(t, err)
	})
}

func TestHandlerUpdate(t *testing.T) {
	v := validator.NewReqValidator()

	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()

		svc := newMockServicer(t)
		svc.On("Update", contextBackground(), "john", member).Return(nil)

		h := NewHandler(svc)
		body := `{"firstName":"John","lastName":"Doe","birthday":"2000-01-01T00:00:00Z"}`
		ctx, rec := echotest.ContextConfig{
			Request:  httptest.NewRequest(http.MethodPut, "/api/v1/members/john", strings.NewReader(body)),
			JSONBody: []byte(body),
			PathValues: echo.PathValues{
				{Name: "username", Value: "john"},
			},
		}.ToContextRecorder(t)
		ctx.Echo().Validator = v

		err := h.update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("bind error - invalid json", func(t *testing.T) {
		svc := newMockServicer(t)
		h := NewHandler(svc)
		ctx, _ := echotest.ContextConfig{
			Request:  httptest.NewRequest(http.MethodPut, "/api/v1/members/john", strings.NewReader("{invalid}")),
			JSONBody: []byte("{invalid}"),
			PathValues: echo.PathValues{
				{Name: "username", Value: "john"},
			},
		}.ToContextRecorder(t)
		ctx.Echo().Validator = v

		err := h.update(ctx)
		assert.Error(t, err)
	})

	t.Run("service error", func(t *testing.T) {
		member, _ := newFixture()

		svc := newMockServicer(t)
		svc.On("Update", contextBackground(), "john", member).Return(ErrorMemberNotFound)

		h := NewHandler(svc)
		body := `{"firstName":"John","lastName":"Doe","birthday":"2000-01-01T00:00:00Z"}`
		ctx, _ := echotest.ContextConfig{
			Request:  httptest.NewRequest(http.MethodPut, "/api/v1/members/john", strings.NewReader(body)),
			JSONBody: []byte(body),
			PathValues: echo.PathValues{
				{Name: "username", Value: "john"},
			},
		}.ToContextRecorder(t)
		ctx.Echo().Validator = v

		err := h.update(ctx)
		assert.Error(t, err)
	})
}

func TestHandlerRemove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := newMockServicer(t)
		svc.On("Remove", contextBackground(), "john").Return(nil)

		h := NewHandler(svc)
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodDelete, "/api/v1/members/john", nil),
			PathValues: echo.PathValues{
				{Name: "username", Value: "john"},
			},
		}.ToContextRecorder(t)

		err := h.remove(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("bind error - invalid body with json content type", func(t *testing.T) {
		svc := newMockServicer(t)
		h := NewHandler(svc)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/members/john", strings.NewReader("{invalid}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx, _ := echotest.ContextConfig{
			Request: req,
		}.ToContextRecorder(t)

		err := h.remove(ctx)
		assert.Error(t, err)
	})

	t.Run("service error", func(t *testing.T) {
		svc := newMockServicer(t)
		svc.On("Remove", contextBackground(), "john").Return(ErrorMemberNotFound)

		h := NewHandler(svc)
		ctx, _ := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodDelete, "/api/v1/members/john", nil),
			PathValues: echo.PathValues{
				{Name: "username", Value: "john"},
			},
		}.ToContextRecorder(t)

		err := h.remove(ctx)
		assert.Error(t, err)
	})
}
