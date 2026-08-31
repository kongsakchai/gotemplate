package authapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/auth"
	pkgconfig "github.com/kongsakchai/gotemplate/pkg/config"
	"github.com/kongsakchai/gotemplate/pkg/serror"
	"github.com/kongsakchai/gotemplate/pkg/validator"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleError(t *testing.T) {
	t.Run("should return NotFound for ErrUserNotFound", func(t *testing.T) {
		appErr := auth.ErrUserNotFound.Err()
		result := handleError(appErr)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.Equal(t, auth.ErrUserNotFound.Code, result.Code)
		assert.Equal(t, auth.ErrUserNotFound.Msg, result.Message)
	})

	t.Run("should return BadRequest for ErrInvalidPassword", func(t *testing.T) {
		appErr := auth.ErrInvalidPassword.Err()
		result := handleError(appErr)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.Equal(t, auth.ErrInvalidPassword.Code, result.Code)
		assert.Equal(t, auth.ErrInvalidPassword.Msg, result.Message)
	})

	t.Run("should return BadRequest for ErrUserAlreadyExists", func(t *testing.T) {
		appErr := auth.ErrUserAlreadyExists.Err()
		result := handleError(appErr)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.Equal(t, auth.ErrUserAlreadyExists.Code, result.Code)
		assert.Equal(t, auth.ErrUserAlreadyExists.Msg, result.Message)
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
	t.Run("should create authApp with dependencies", func(t *testing.T) {
		a := NewApp(Deps{})
		assert.NotNil(t, a)
	})
}

func TestAuthAppRouteRegistration(t *testing.T) {
	t.Run("should register routes without panic", func(t *testing.T) {
		cfg := pkgconfig.Config{}
		echoApp := app.NewEchoApp(cfg)
		a := NewApp(Deps{})

		assert.NotPanics(t, func() {
			a.RegisterRoute(echoApp)
		})
	})
}

func TestAuthAppHandlers(t *testing.T) {
	// newContext creates an echotest context with a JSON body and registers
	// a validator so app.Request's Validate step works as in production.
	newContext := func(body string) (*echo.Context, *httptest.ResponseRecorder) {
		ctx, rec := echotest.ContextConfig{
			JSONBody: []byte(body),
		}.ToContextRecorder(t)
		ctx.Echo().Validator = validator.NewValidator()
		return ctx, rec
	}

	newApp := func(sv *mockServicer) *authApp {
		a := NewApp(Deps{})
		a.sv = sv
		return a
	}

	// === register handler ===

	t.Run("register - success returns 201 Created", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().Register(mock.Anything, "testuser", "testpass").Return(nil)

		ctx, rec := newContext(`{"username":"testuser","password":"testpass"}`)
		err := newApp(sv).register(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp app.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
	})

	t.Run("register - invalid JSON returns BadRequest", func(t *testing.T) {
		ctx, _ := newContext(`{invalid json}`)
		err := newApp(newMockServicer(t)).register(ctx)

		var appErr app.Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
	})

	t.Run("register - missing required fields returns BadRequest", func(t *testing.T) {
		ctx, _ := newContext(`{}`)
		err := newApp(newMockServicer(t)).register(ctx)

		var appErr app.Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
	})

	t.Run("register - service error returns serror", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().Register(mock.Anything, "existing", "secret").Return(auth.ErrUserAlreadyExists.Err())

		ctx, _ := newContext(`{"username":"existing","password":"secret"}`)
		err := newApp(sv).register(ctx)

		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, auth.ErrUserAlreadyExists.Code, serr.Code())
		assert.Equal(t, auth.ErrUserAlreadyExists.Msg, serr.Msg())
	})

	// === login handler ===

	t.Run("login - success returns 200 with token", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().Login(mock.Anything, "testuser", "testpass").Return("jwt-token-value", nil)

		ctx, rec := newContext(`{"username":"testuser","password":"testpass"}`)
		err := newApp(sv).login(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Code    string `json:"code"`
			Success bool   `json:"success"`
			Data    struct {
				Token string `json:"token"`
			} `json:"data"`
		}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
		assert.Equal(t, "jwt-token-value", resp.Data.Token)
	})

	t.Run("login - invalid JSON returns bind error", func(t *testing.T) {
		ctx, _ := newContext(`bad json`)
		err := newApp(newMockServicer(t)).login(ctx)

		assert.Error(t, err)
	})

	t.Run("login - service error returns serror", func(t *testing.T) {
		sv := newMockServicer(t)
		sv.EXPECT().Login(mock.Anything, "ghost", "secret").Return("", auth.ErrUserNotFound.Err())

		ctx, _ := newContext(`{"username":"ghost","password":"secret"}`)
		err := newApp(sv).login(ctx)

		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, auth.ErrUserNotFound.Code, serr.Code())
		assert.Equal(t, auth.ErrUserNotFound.Msg, serr.Msg())
	})
}
