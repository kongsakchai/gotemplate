package mockutil

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type Timer struct {
	mock.Mock
}

func NewTimer(t *testing.T) *Timer {
	mock := &Timer{}
	mock.Mock.Test(t)

	return mock
}
