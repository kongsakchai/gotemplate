package app

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kongsakchai/gotemplate/pkg/token"
	"github.com/labstack/echo/v5"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

func AuthMiddleware(verifier token.Verifier) echo.MiddlewareFunc {
	prefixLen := len(bearerPrefix)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			values := ctx.Request().Header.Values(authorizationHeader)
			if len(values) == 0 {
				return Unauthorized(MissingTokenCode, MissingTokenMsg, nil)
			}

			for _, v := range values {
				claims, err := verifier.Verify(v[prefixLen:])
				if errors.Is(err, jwt.ErrTokenExpired) {
					return Unauthorized(TokenExpiredCode, TokenExpiredMsg, err)
				}
				if err != nil {
					return Unauthorized(UnauthorizedCode, UnauthorizedMsg, err)
				}

				for key, value := range claims {
					ctx.Set(key, value)
				}
			}

			return next(ctx)
		}
	}
}
