package member

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func setupStorage(t *testing.T) *storage {
	t.Helper()

	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE member (
		username TEXT PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		birthday datetime,
		register_date datetime
	)`)
	require.NoError(t, err)

	return NewStorage(db)
}

func TestStorageCreate(t *testing.T) {
	t.Run("should create member successfully", func(t *testing.T) {
		s := setupStorage(t)

		m := Member{
			Username: "newuser",
		}
		err := s.Create(t.Context(), m)
		assert.NoError(t, err)
	})
}

func TestStorageUpdate(t *testing.T) {
	t.Run("should update existing member", func(t *testing.T) {
		s := setupStorage(t)

		_, err := s.db.ExecContext(t.Context(),
			"INSERT INTO member (username, first_name, last_name, birthday, register_date) VALUES (?, ?, ?, ?, ?)",
			"toupdate", "Old", "Name", "2000-01-01T00:00:00Z", "2025-01-01T00:00:00Z")
		require.NoError(t, err)

		m := Member{
			Username:     "toupdate",
			FirstName:    "Updated",
			LastName:     "Name",
			Birthday:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			RegisterDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		err = s.Update(t.Context(), m)
		assert.NoError(t, err)
	})
}

func TestStorageRemove(t *testing.T) {
	t.Run("should remove existing member", func(t *testing.T) {
		s := setupStorage(t)

		_, err := s.db.ExecContext(t.Context(),
			"INSERT INTO member (username) VALUES (?)", "todelete")
		require.NoError(t, err)

		err = s.Remove(t.Context(), "todelete")
		assert.NoError(t, err)
	})
}

func TestNewStorage(t *testing.T) {
	t.Run("should create storage", func(t *testing.T) {
		s := setupStorage(t)
		assert.NotNil(t, s)
	})
}

func TestStorageMembers(t *testing.T) {
	t.Run("should return members", func(t *testing.T) {
		s := setupStorage(t)

		_, err := s.db.ExecContext(t.Context(),
			"INSERT INTO member (username, first_name, last_name, birthday, register_date) VALUES (?, ?, ?, ?, ?)",
			"toupdate", "Old", "Name", "2000-01-01T00:00:00Z", "2025-01-01T00:00:00Z")
		require.NoError(t, err)

		members, err := s.Members(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(members))
	})
}

func TestStorageMember(t *testing.T) {
	t.Run("should return member", func(t *testing.T) {
		s := setupStorage(t)

		_, err := s.db.ExecContext(t.Context(),
			"INSERT INTO member (username, first_name, last_name, birthday, register_date) VALUES (?, ?, ?, ?, ?)",
			"test", "Old", "Name", "2000-01-01T00:00:00Z", "2025-01-01T00:00:00Z")
		require.NoError(t, err)

		member, found, err := s.Member(t.Context(), "test")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "Old", member.FirstName)
	})

	t.Run("should return empty when no data", func(t *testing.T) {
		s := setupStorage(t)

		member, found, err := s.Member(t.Context(), "test")
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Equal(t, "", member.FirstName)
	})
}
