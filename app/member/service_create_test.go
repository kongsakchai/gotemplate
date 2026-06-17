package member

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		member, now := newFixture()
		expected := member
		expected.RegisterDate = now

		svc := newServiceWithMocks(t, func(c *mockClock) {
			c.On("Now").Return(now)
		}, func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(Member{}, false, nil)
			m.EXPECT().Create(contextBackground(), expected).Return(nil)
		})

		err := svc.Create(contextBackground(), member)
		assert.NoError(t, err)
	})

	t.Run("age too young", func(t *testing.T) {
		member, now := newFixture()
		member.Birthday = now.AddDate(0, 0, -1)

		svc := newServiceWithMocks(t, func(c *mockClock) {
			c.On("Now").Return(now)
		}, noStorage())

		err := svc.Create(contextBackground(), member)
		assert.ErrorIs(t, err, ErrorMinAge)
	})

	t.Run("age too old", func(t *testing.T) {
		member, now := newFixture()
		member.Birthday = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

		svc := newServiceWithMocks(t, func(c *mockClock) {
			c.On("Now").Return(now)
		}, noStorage())

		err := svc.Create(contextBackground(), member)
		assert.ErrorIs(t, err, ErrorMaxAge)
	})

	t.Run("duplicate username", func(t *testing.T) {
		member, now := newFixture()

		svc := newServiceWithMocks(t, func(c *mockClock) {
			c.On("Now").Return(now)
		}, func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(Member{}, true, nil)
		})

		err := svc.Create(contextBackground(), member)
		assert.ErrorIs(t, err, ErrorDuplicate)
	})

	t.Run("storage error on member check", func(t *testing.T) {
		member, now := newFixture()

		svc := newServiceWithMocks(t, func(c *mockClock) {
			c.On("Now").Return(now)
		}, func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(Member{}, false, errors.New("db err"))
		})

		err := svc.Create(contextBackground(), member)
		assert.ErrorContains(t, err, "db err")
	})

	t.Run("storage error on create", func(t *testing.T) {
		member, now := newFixture()
		expected := member
		expected.RegisterDate = now

		svc := newServiceWithMocks(t, func(c *mockClock) {
			c.On("Now").Return(now)
		}, func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(Member{}, false, nil)
			m.EXPECT().Create(contextBackground(), expected).Return(errors.New("insert err"))
		})

		err := svc.Create(contextBackground(), member)
		assert.ErrorContains(t, err, "insert err")
	})
}
