package auth

import (
	"context"
	"time"

	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
	"github.com/kongsakchai/gotemplate/pkg/serror"
)

var (
	ErrUserAlreadyExists = serror.NewCoded("2000", "user already exists")
	ErrUserNotFound      = serror.NewCoded("2001", "user not found")
	ErrInvalidPassword   = serror.NewCoded("2002", "invalid password")
)

type User struct {
	Id        string
	Username  string
	Password  string
	CreatedAt time.Time
}

//mockery:generate: true
type Storager interface {
	CreateUser(ctx context.Context, user User) error
	FindUserByUsername(ctx context.Context, username string) (User, bool, error)
}

//mockery:generate: true
type Hasher = hash.Hasher

//mockery:generate: true
type Signer = jwttoken.Signer

type Servicer interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) error
}
