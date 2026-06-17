package member

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(member, true, nil)
		})

		got, err := svc.Member(contextBackground(), "john")
		assert.NoError(t, err)
		assert.Equal(t, member, got)
	})

	t.Run("not found", func(t *testing.T) {
		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "unknown").Return(Member{}, false, nil)
		})

		_, err := svc.Member(contextBackground(), "unknown")
		assert.ErrorIs(t, err, ErrorMemberNotFound)
	})

	t.Run("storage error", func(t *testing.T) {
		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Member(contextBackground(), "john").Return(Member{}, false, errors.New("db err"))
		})

		_, err := svc.Member(contextBackground(), "john")
		assert.ErrorContains(t, err, "db err")
	})
}
