package jwttoken

import (
	"errors"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kongsakchai/gotemplate/pkg/config"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token is expired")
	ErrInvalidClaims = errors.New("invalid claims")
	ErrSubIsRequired = errors.New("sub is required")
	ErrJTIIsRequired = errors.New("jti is required")
)

type Signer interface {
	Sign(sub, id string, extra map[string]any) (string, error)
}

type Verifier interface {
	Verify(tokenString string) (map[string]any, error)
}

type jwtToken struct {
	method       jwt.SigningMethod
	secretKey    []byte
	verifyKey    []byte
	expired      time.Duration
	issuer       string
	audience     string
	parseOptions []jwt.ParserOption
}

func NewJWTToken(cfg config.JWT) *jwtToken {
	if cfg.Method == "" {
		panic("Missing jwt method")
	}
	if cfg.SecretKey == "" {
		panic("Missing jwt secret key")
	}
	if cfg.VerifyKey == "" {
		panic("Missing jwt verify key")
	}
	parseOptions := []jwt.ParserOption{
		jwt.WithValidMethods([]string{cfg.Method}),
		jwt.WithExpirationRequired(),
	}
	if cfg.Audience != "" {
		parseOptions = append(parseOptions, jwt.WithAudience(cfg.Audience))
	}
	if cfg.Issuer != "" {
		parseOptions = append(parseOptions, jwt.WithIssuer(cfg.Issuer))
	}

	return &jwtToken{
		method:    jwt.GetSigningMethod(cfg.Method),
		secretKey: []byte(cfg.SecretKey),
		verifyKey: []byte(cfg.VerifyKey),
		expired:   cfg.Expired,
		issuer:    cfg.Issuer,
		audience:  cfg.Audience,
	}
}

func (s *jwtToken) Sign(sub, jti string, extra map[string]any) (string, error) {
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
		return "", err
	}
	return tokenString, nil
}

func validateKey(key []byte) func(token *jwt.Token) (any, error) {
	return func(token *jwt.Token) (any, error) {
		return key, nil
	}
}

func (v *jwtToken) Verify(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString,
		validateKey(v.verifyKey),
		v.parseOptions...,
	)
	if err == jwt.ErrTokenExpired {
		return nil, ErrTokenExpired
	}
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	return claims, nil
}
