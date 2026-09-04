package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/serror"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestService_Register(t *testing.T) {
	t.Run("should register user successfully", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "testuser"
		password := "password123"

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(User{}, false, nil)
		mockHasher.EXPECT().HashPassword(password).Return("$2a$hash", nil)
		mockStorager.EXPECT().CreateUser(ctx, mock.MatchedBy(func(u User) bool {
			return u.Username == username && u.Password == "$2a$hash"
		})).Return(nil)

		err := svc.Register(ctx, username, password)
		assert.NoError(t, err)
	})

	t.Run("should return ErrUserAlreadyExists when user already exists", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "existinguser"

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(User{Id: "1", Username: username}, true, nil)

		err := svc.Register(ctx, username, "password123")
		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, serr.Code(), ErrUserAlreadyExists.Code)
		assert.Equal(t, serr.Msg(), ErrUserAlreadyExists.Msg)
	})

	t.Run("should propagate storage error from FindUserByUsername", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "testuser"
		dbErr := errors.New("database error")

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(User{}, false, dbErr)

		err := svc.Register(ctx, username, "password123")
		assert.Same(t, dbErr, err)
	})

	t.Run("should propagate hasher error", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "testuser"
		hashErr := errors.New("hash error")

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(User{}, false, nil)
		mockHasher.EXPECT().HashPassword("password123").Return("", hashErr)

		err := svc.Register(ctx, username, "password123")
		assert.Same(t, hashErr, err)
	})
}

func TestService_Login(t *testing.T) {
	t.Run("should login and return token successfully", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)
		mockSigner := newMockSigner(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
			Signer:   mockSigner,
		})

		ctx := context.Background()
		username := "testuser"
		password := "password123"
		user := User{Id: "user-1", Username: username, Password: "$2a$hashed"}

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(user, true, nil)
		mockHasher.EXPECT().ComparePassword(password, "$2a$hashed").Return(nil)
		mockSigner.EXPECT().Sign("user-1", mock.MatchedBy(func(id string) bool { return id != "" }), map[string]any(nil)).Return("token123", nil)

		token, err := svc.Login(ctx, username, password)
		assert.NoError(t, err)
		assert.Equal(t, "token123", token)
	})

	t.Run("should return ErrUserNotFound when user does not exist", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "unknownuser"

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(User{}, false, nil)

		token, err := svc.Login(ctx, username, "password123")
		assert.Empty(t, token)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, serr.Code(), ErrUserNotFound.Code)
		assert.Equal(t, serr.Msg(), ErrUserNotFound.Msg)
	})

	t.Run("should return ErrInvalidPassword when password is wrong", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "testuser"
		user := User{Id: "user-1", Username: username, Password: "$2a$hashed"}

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(user, true, nil)
		mockHasher.EXPECT().ComparePassword("wrongpass", "$2a$hashed").Return(hash.ErrMismatched)

		token, err := svc.Login(ctx, username, "wrongpass")
		assert.Empty(t, token)
		serr, ok := serror.As(err)
		assert.True(t, ok)
		assert.Equal(t, serr.Code(), ErrInvalidPassword.Code)
		assert.Equal(t, serr.Msg(), ErrInvalidPassword.Msg)
	})

	t.Run("should propagate hasher error that is not ErrMismatched", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "testuser"
		user := User{Id: "user-1", Username: username, Password: "$2a$hashed"}
		hasherErr := errors.New("hash comparison error")

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(user, true, nil)
		mockHasher.EXPECT().ComparePassword("password123", "$2a$hashed").Return(hasherErr)

		token, err := svc.Login(ctx, username, "password123")
		assert.Empty(t, token)
		assert.Same(t, hasherErr, err)
	})

	t.Run("should propagate storage error from FindUserByUsername", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
		})

		ctx := context.Background()
		username := "testuser"
		dbErr := errors.New("database error")

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(User{}, false, dbErr)

		token, err := svc.Login(ctx, username, "password123")
		assert.Empty(t, token)
		assert.Same(t, dbErr, err)
	})

	t.Run("should propagate signer error", func(t *testing.T) {
		mockStorager := newMockStorager(t)
		mockHasher := newMockHasher(t)
		mockSigner := newMockSigner(t)

		svc := NewService(ServiceDeps{
			Storager: mockStorager,
			Hasher:   mockHasher,
			Signer:   mockSigner,
		})

		ctx := context.Background()
		username := "testuser"
		user := User{Id: "user-1", Username: username, Password: "$2a$hashed"}
		signerErr := errors.New("signing error")

		mockStorager.EXPECT().FindUserByUsername(ctx, username).Return(user, true, nil)
		mockHasher.EXPECT().ComparePassword("password123", "$2a$hashed").Return(nil)
		mockSigner.EXPECT().Sign("user-1", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return("", signerErr)

		token, err := svc.Login(ctx, username, "password123")
		assert.Empty(t, token)
		assert.Same(t, signerErr, err)
	})
}
