package member

import (
	"context"
	"errors"
	"time"
)

var (
	ErrorMinAge         = errors.New("min age limit")
	ErrorMaxAge         = errors.New("max age limit")
	ErrorDuplicate      = errors.New("duplicate username")
	ErrorMemberNotFound = errors.New("member not found")
)

type Member struct {
	Username     string    `json:"username"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Birthday     time.Time `json:"birthday"`
	RegisterDate time.Time `json:"registerDate"`
}

//mockery:generate: true
type Storager interface {
	Members(ctx context.Context) ([]Member, error)
	Member(ctx context.Context, username string) (Member, bool, error)
	Create(ctx context.Context, member Member) error
	Remove(ctx context.Context, username string) error
	Update(ctx context.Context, member Member) error
}

//mockery:generate: true
type Servicer interface {
	Members(ctx context.Context) ([]Member, error)
	Member(ctx context.Context, username string) (Member, error)
	Create(ctx context.Context, member Member) error
	Remove(ctx context.Context, username string) error
	Update(ctx context.Context, username string, member Member) error
}

//mockery:generate: true
type Clock interface {
	Now() time.Time
}
