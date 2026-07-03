package member

import (
	"context"
	"time"

	"github.com/kongsakchai/gotemplate/domain"
)

var (
	MemberNotFoundCode      = "2000"
	InvalidAgeCode          = "2001"
	UsernameUnavailableCode = "2002"

	ErrMemberNotFound      = domain.NewError(MemberNotFoundCode, "member not found")
	ErrInvalidAge          = domain.NewError(InvalidAgeCode, "age must be between 15 and 60")
	ErrUsernameUnavailable = domain.NewError(UsernameUnavailableCode, "username unavailable")
)

type Member struct {
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateMemberPayload struct {
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthday  time.Time `json:"birthday"`
}

type UpdateMemberPayload struct {
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthday  time.Time `json:"birthday"`
}

//mockery:generate: true
type Storager interface {
	GetMembers(ctx context.Context) ([]Member, error)
	GetMember(ctx context.Context, username string) (Member, bool, error)
	CreateMember(ctx context.Context, member CreateMemberPayload) error
	RemoveMember(ctx context.Context, username string) error
	UpdateMember(ctx context.Context, member UpdateMemberPayload) error
}

//mockery:generate: true
type Servicer interface {
	GetMembers(ctx context.Context) ([]Member, error)
	GetMember(ctx context.Context, username string) (Member, error)
	CreateMember(ctx context.Context, member CreateMemberPayload) error
	RemoveMember(ctx context.Context, username string) error
	UpdateMember(ctx context.Context, member UpdateMemberPayload) error
}
