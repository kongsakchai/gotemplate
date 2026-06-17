package member

import (
	"context"
	"time"

	mock "github.com/stretchr/testify/mock"
)

func newFixture() (Member, time.Time) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return Member{
		Username:  "john",
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  birthday,
	}, now
}

type mockClockFn func(*mockClock)
type mockStorageFn func(*mockStorager)

func newServiceWithMocks(t interface {
	mock.TestingT
	Cleanup(func())
}, clockFn mockClockFn, storageFn mockStorageFn) *service {
	clock := &mockClock{}
	clock.Test(t)

	storage := newMockStorager(t)

	if clockFn != nil {
		clockFn(clock)
	}
	if storageFn != nil {
		storageFn(storage)
	}

	return NewService(storage, clock)
}

func noClock() mockClockFn     { return nil }
func noStorage() mockStorageFn { return nil }

func contextBackground() context.Context {
	return context.Background()
}
