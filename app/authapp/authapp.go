package authapp

import (
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/auth"
	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
	"github.com/labstack/echo/v5"
)

//mockery:generate: true
type Servicer = auth.Servicer

type authApp struct {
	sv Servicer
}

type Deps struct {
	DB     *sqlx.DB
	Hasher hash.Hasher
	Signer jwttoken.Signer
}

func NewApp(deps Deps) *authApp {
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

func (a *authApp) RegisterRoute(echo *app.EchoApp) {
	g := echo.Group("/api/v1/auth", app.ErrorMiddleware(handleError))
	app.POST(g, "auth_register", "/register", a.register)
	app.POST(g, "auth_login", "/login", a.login)
}

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *authApp) register(c *echo.Context) error {
	req, err := app.Request[registerRequest](c)
	if err != nil {
		return err
	}
	if err := a.sv.Register(c.Request().Context(), req.Username, req.Password); err != nil {
		return err
	}
	return app.Created(c, nil)
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (a *authApp) login(c *echo.Context) error {
	req, err := app.Request[loginRequest](c)
	if err != nil {
		return err
	}
	token, err := a.sv.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return err
	}

	return app.Ok(c, loginResponse{Token: token})
}
