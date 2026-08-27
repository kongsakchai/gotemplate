package jwttoken

import (
	"errors"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kongsakchai/gotemplate/pkg/errs"
)

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

type JWTOptions struct {
	Secret    string
	VerifyKey string
	Method    string
	Expired   time.Duration
	Issuer    string
	Audience  string
}

type jwtManager struct {
	method    jwt.SigningMethod
	secretKey []byte
	verifyKey string
	expired   time.Duration
	issuer    string
	audience  string
}

func NewJWTManager(opt JWTOptions) *jwtManager {
	return &jwtManager{
		method:    jwt.GetSigningMethod(opt.Method),
		secretKey: []byte(opt.Secret),
		verifyKey: opt.VerifyKey,
		expired:   opt.Expired,
		issuer:    opt.Issuer,
		audience:  opt.Audience,
	}
}

func (s *jwtManager) Sign(sub, jti string, extra map[string]any) (string, error) {
	if sub == "" {
		return "", ErrSubIsRequired
	}
	if jti == "" {
		return "", ErrJTIIsRequired
	}

	now := time.Now()
	payload := make(map[string]any)
	maps.Copy(payload, extra)

	payload["sub"] = sub
	payload["jti"] = jti
	payload["exp"] = now.Add(s.expired).Unix()
	payload["iat"] = now.Unix()
	payload["iss"] = s.issuer
	payload["aud"] = s.audience

	token := jwt.NewWithClaims(s.method, jwt.MapClaims(payload))
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", errs.From(err)
	}
	return tokenString, nil
}

func validateKey(key string) func(token *jwt.Token) (any, error) {
	return func(token *jwt.Token) (any, error) {
		return key, nil
	}
}

func (v *jwtManager) Verify(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString,
		validateKey(v.verifyKey),
		jwt.WithValidMethods([]string{v.method.Alg()}),
		jwt.WithAudience(v.audience),
		jwt.WithIssuer(v.issuer),
		jwt.WithExpirationRequired(),
	)
	if err == jwt.ErrTokenExpired {
		return nil, ErrTokenExpired
	}
	if err != nil {
		return nil, errs.From(err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errs.New("invalid claims")
	}
	return claims, nil
}
