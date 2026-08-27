package hash

import (
	"errors"

	"github.com/kongsakchai/gotemplate/pkg/serror"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrMismatched = errors.New("comparison failed")
)

type Hasher interface {
	HashPassword(pwd string) (string, error)
	ComparePassword(pwd, hash string) error
}

type hasher struct {
	cost int
}

func NewHasher(cost int) Hasher {
	return &hasher{cost: cost}
}

func (h *hasher) HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), h.cost)
	if err != nil {
		return "", serror.From(err)
	}
	return string(hash), nil
}

func (h *hasher) ComparePassword(pwd, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return serror.From(ErrMismatched)
	}
	return serror.From(err)
}
