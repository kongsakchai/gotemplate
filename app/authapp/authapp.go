package authapp

import (
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/auth"
	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
	"github.com/kongsakchai/gotemplate/pkg/serror"
	"github.com/labstack/echo/v5"
)

type authApp struct {
	sv auth.Service
}

type Deps struct {
	DB     *sqlx.DB
	Hasher hash.Hasher
	Signer jwttoken.Signer
}

func NewAuthApp(deps Deps) *authApp {
	st := auth.NewStorage(deps.DB)
	sv := auth.NewService(auth.ServiceDeps{
		Storager: st,
		Hasher:   deps.Hasher,
		Signer:   deps.Signer,
	})

	return &authApp{
		sv: sv,
	}
}

func (a *authApp) RegisterRoute(app *app.EchoApp) {
	app.POST("auth_register", "/api/v1/auth/register", a.register)
	app.POST("auth_login", "/api/v1/auth/login", a.login)
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *authApp) register(c *echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := a.sv.Register(c.Request().Context(), req.Username, req.Password); err != nil {
		return handleError(err)
	}
	return app.Created(c, nil)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (a *authApp) login(c *echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	token, err := a.sv.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return handleError(err)
	}

	return app.Ok(c, loginResponse{Token: token})
}

func handleError(err error) error {
	if e, ok := serror.As(err); ok {
		switch e.Code() {
		case auth.ErrUserNotFound.Code:
			return app.NotFound(e.Code(), e.Message(), e)
		case auth.ErrInvalidPassword.Code, auth.ErrUserAlreadyExists.Code:
			return app.BadRequest(e.Code(), e.Message(), e)
		}
	}

	return app.InternalError(app.InternalErrorCode, app.InternalErrorMsg, err)
}
