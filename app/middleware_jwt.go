package app

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
	"github.com/labstack/echo/v5"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

func AuthMiddleware(verifier jwttoken.Verifier) echo.MiddlewareFunc {
	prefixLen := len(bearerPrefix)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			values := ctx.Request().Header.Values(authorizationHeader)
			if len(values) == 0 {
				return Unauthorized(MissingTokenCode, MissingTokenMsg, nil)
			}

			for _, v := range values {
				if len(v) <= prefixLen {
					return Unauthorized(InvalidTokenCode, InvalidTokenMsg, nil)
				}
				claims, err := verifier.Verify(v[prefixLen:])
				if errors.Is(err, jwt.ErrTokenExpired) {
					return Unauthorized(TokenExpiredCode, TokenExpiredMsg, err)
				}
				if err != nil {
					return Unauthorized(UnauthorizedCode, UnauthorizedMsg, err)
				}

				reqCtx := ctx.Request().Context()
				for key, value := range claims {
					ctx.Set(key, value)
					reqCtx = context.WithValue(reqCtx, key, value)
				}
				ctx.SetRequest(ctx.Request().WithContext(reqCtx))
			}

			return next(ctx)
		}
	}
}
