package token

import "errors"

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token is expired")
	ErrClaimsInvalid = errors.New("claims are invalid")
	ErrSubIsRequired = errors.New("sub is required")
	ErrJTIIsRequired = errors.New("jti is required")
)

type Signer interface {
	Sign(sub, jti string, extra map[string]any) (string, error)
}

type Verifier interface {
	Verify(tokenString string) (map[string]any, error)
}
