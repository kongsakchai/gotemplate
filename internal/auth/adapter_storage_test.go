package auth

import (
	"context"
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/serror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	t.Run("should create storage with db", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		assert.NotNil(t, st)
		assert.NotNil(t, st.db)
	})
}

func TestStorage_FindUserByUsername(t *testing.T) {
	ctx := context.Background()

	t.Run("should return user when found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at"}).
			AddRow("user-1", "testuser", "$2a$hash", now)

		mock.ExpectQuery("SELECT id, username, password_hash, created_at FROM user WHERE username").
			WithArgs("testuser").
			WillReturnRows(rows)

		user, found, err := st.FindUserByUsername(ctx, "testuser")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "user-1", user.Id)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "$2a$hash", user.Password)
	})

	t.Run("should return false when user not found (ErrNoRows)", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT id, username, password_hash, created_at FROM user WHERE username").
			WithArgs("unknownuser").
			WillReturnError(sql.ErrNoRows)

		user, found, err := st.FindUserByUsername(ctx, "unknownuser")
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Equal(t, User{}, user)
	})

	t.Run("should wrap and return database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectQuery("SELECT id, username, password_hash, created_at FROM user WHERE username").
			WithArgs("testuser").
			WillReturnError(assert.AnError)

		user, found, err := st.FindUserByUsername(ctx, "testuser")
		assert.Error(t, err)
		assert.False(t, found)
		assert.Equal(t, User{}, user)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		_ = serr
	})
}

func TestStorage_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("INSERT INTO user").
			WithArgs("user-1", "testuser", "$2a$hash").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = st.CreateUser(ctx, User{
			Id:       "user-1",
			Username: "testuser",
			Password: "$2a$hash",
		})
		assert.NoError(t, err)
	})

	t.Run("should wrap and return database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		st := NewStorage(sqlxDB)

		mock.ExpectExec("INSERT INTO user").
			WithArgs("user-1", "testuser", "$2a$hash").
			WillReturnError(assert.AnError)

		err = st.CreateUser(ctx, User{
			Id:       "user-1",
			Username: "testuser",
			Password: "$2a$hash",
		})
		assert.NotNil(t, err)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		_ = serr
	})
}

func TestUserRecord_toUser(t *testing.T) {
	t.Run("should convert userRecord to User", func(t *testing.T) {
		now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		record := userRecord{
			Id:        "user-1",
			Username:  "testuser",
			Password:  "$2a$hashed",
			CreatedAt: now,
		}

		user := record.toUser()

		assert.Equal(t, "user-1", user.Id)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "$2a$hashed", user.Password)
		assert.Equal(t, now, user.CreatedAt)
	})
}
