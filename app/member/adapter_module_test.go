package member

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClock2 struct{}

func (m *mockClock2) Now() time.Time { return time.Time{} }

func TestNewModule(t *testing.T) {
	t.Run("should create module with handler", func(t *testing.T) {
		db, err := sqlx.Open("sqlite", ":memory:")
		require.NoError(t, err)
		t.Cleanup(func() { db.Close() })

		mod := NewModule(External{
			DB:    db,
			Clock: &mockClock2{},
		})
		assert.NotNil(t, mod)
		assert.NotNil(t, mod.Handler)
	})
}
