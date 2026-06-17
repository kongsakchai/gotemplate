package member

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(member, true, nil)
			m.EXPECT().Update(contextBackground(), member).Return(nil)
		})

		err := svc.Update(contextBackground(), "john", member)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		member, _ := newFixture()

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "unknown").Return(Member{}, false, nil)
		})

		err := svc.Update(contextBackground(), "unknown", member)
		assert.ErrorIs(t, err, ErrorMemberNotFound)
	})

	t.Run("storage error on member check", func(t *testing.T) {
		member, _ := newFixture()

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(Member{}, false, errors.New("db err"))
		})

		err := svc.Update(contextBackground(), "john", member)
		assert.ErrorContains(t, err, "db err")
	})

	t.Run("storage error on update", func(t *testing.T) {
		member, _ := newFixture()

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(member, true, nil)
			m.EXPECT().Update(contextBackground(), member).Return(errors.New("update err"))
		})

		err := svc.Update(contextBackground(), "john", member)
		assert.ErrorContains(t, err, "update err")
	})
}
