package member

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceMembers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		member, _ := newFixture()
		expected := []Member{member}

		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Members(contextBackground()).Return(expected, nil)
		})

		got, err := svc.Members(contextBackground())
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("storage error", func(t *testing.T) {
		svc := newServiceWithMocks(t, noClock(), func(m *mockStorager) {
			m.EXPECT().Members(contextBackground()).Return(nil, errors.New("db err"))
		})

		_, err := svc.Members(contextBackground())
		assert.ErrorContains(t, err, "db err")
	})
}
