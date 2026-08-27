package auth

import (
	"context"

	"github.com/kongsakchai/gotemplate/pkg/serror"
)

var (
	ErrUserAlreadyExists = serror.NewCoded("2000", "user already exists")
	ErrUserNotFound      = serror.NewCoded("2001", "user not found")
	ErrInvalidPassword   = serror.NewCoded("2002", "invalid password")
)

type User struct {
	Id       string
	Username string
	Password string
}

type Storager interface {
	CreateUser(ctx context.Context, user User) error
	FindUserByUsername(ctx context.Context, username string) (User, bool, error)
}

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) error
}
