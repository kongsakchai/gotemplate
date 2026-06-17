package generate

import "github.com/google/uuid"

type uuidGenerator struct{}

func NewUUID() *uuidGenerator {
	return &uuidGenerator{}
}

func (u *uuidGenerator) GenUUID() string {
	return uuid.NewString()
}
