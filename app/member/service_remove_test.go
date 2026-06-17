package member

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceRemove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(member, true, nil)
			m.EXPECT().Remove(contextBackground(), "john").Return(nil)
		})

		err := svc.Remove(contextBackground(), "john")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "unknown").Return(Member{}, false, nil)
		})

		err := svc.Remove(contextBackground(), "unknown")
		assert.ErrorIs(t, err, ErrorMemberNotFound)
	})

	t.Run("storage error on member check", func(t *testing.T) {
		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(Member{}, false, errors.New("db err"))
		})

		err := svc.Remove(contextBackground(), "john")
		assert.ErrorContains(t, err, "db err")
	})

	t.Run("storage error on remove", func(t *testing.T) {
		member, _ := newFixture()

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(member, true, nil)
			m.EXPECT().Remove(contextBackground(), "john").Return(errors.New("delete err"))
		})

		err := svc.Remove(contextBackground(), "john")
		assert.ErrorContains(t, err, "delete err")
	})
}
