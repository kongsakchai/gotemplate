package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
)

type service struct {
	st Storager

	hasher Hasher
	signer Signer
}

type ServiceDeps struct {
	Storager Storager
	Hasher   hash.Hasher
	Signer   jwttoken.Signer
}

func NewService(deps ServiceDeps) *service {
	return &service{
		st:     deps.Storager,
		hasher: deps.Hasher,
		signer: deps.Signer,
	}
}

func (s *service) Register(ctx context.Context, username, password string) error {
	_, found, err := s.st.FindUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if found {
		return ErrUserAlreadyExists.Err()
	}
	hashedPassword, err := s.hasher.HashPassword(password)
	if err != nil {
		return err
	}
	return s.st.CreateUser(ctx, User{
		Id:       uuid.NewString(),
		Username: username,
		Password: hashedPassword,
	})
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	user, found, err := s.st.FindUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrUserNotFound.Err()
	}
	if err := s.hasher.ComparePassword(password, user.Password); err != nil {
		if errors.Is(err, hash.ErrMismatched) {
			return "", ErrInvalidPassword.Err()
		}
		return "", err
	}

	token, err := s.signer.Sign(user.Id, uuid.NewString(), nil)
	if err != nil {
		return "", err
	}

	return token, nil
}
