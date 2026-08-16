package password

import (
	"errors"

	"github.com/kongsakchai/gotemplate/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type Hasher interface {
	Hash(pwd string) (string, error)
	Compare(pwd, hash string) error
}

type hasher struct {
	cost int
}

func NewHasher(cost int) Hasher {
	return &hasher{cost: cost}
}

func (h *hasher) Hash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), h.cost)
	if err != nil {
		return "", errs.From(err)
	}
	return string(hash), nil
}

func (h *hasher) Compare(pwd, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidPassword
	}
	return errs.From(err)
}
