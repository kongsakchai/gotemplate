package mockutil

import (
	"testing"
	"time"

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

func (m *Timer) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
