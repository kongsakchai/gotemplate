package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserByName(t *testing.T) {
	t.Run("should return user when user exists", func(t *testing.T) {
		storage := NewStorage()
		user, err := storage.UserByName("john")

		assert.True(t, err.IsEmpty())
		assert.Equal(t, "John", user.Name)
		assert.Equal(t, 30, user.Age)
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {
		storage := NewStorage()
		user, err := storage.UserByName("jane")

		assert.Equal(t, "4001", err.Code)
		assert.Equal(t, "user not found", err.Message)
		assert.Empty(t, user)
	})
}

func TestUsers(t *testing.T) {
	t.Run("should return list of users", func(t *testing.T) {
		storage := NewStorage()
		users, err := storage.Users()

		assert.True(t, err.IsEmpty())
		assert.Len(t, users, 1)
		assert.Equal(t, "John", users[0].Name)
		assert.Equal(t, 30, users[0].Age)
	})
}
