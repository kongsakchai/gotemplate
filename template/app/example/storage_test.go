package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	s := NewStorage()

	assert.NotNil(t, s)
	assert.NotNil(t, s.users)
	assert.Len(t, s.users, 0)
}

func TestStorageCreateUser(t *testing.T) {
	s := NewStorage()
	user := User{FirstName: "john", LastName: "doe", Age: 30}

	err := s.CreateUser(user)

	assert.NoError(t, err)
	assert.Len(t, s.users, 1)
	assert.Equal(t, user, s.users[0])
}

func TestStorageUserByName(t *testing.T) {
	t.Run("should return user when name exists", func(t *testing.T) {
		s := NewStorage()
		expected := User{FirstName: "john", LastName: "doe", Age: 30}
		s.users = append(s.users, expected)

		user, err := s.UserByName("john")

		assert.NoError(t, err)
		assert.Equal(t, expected, user)
	})

	t.Run("should return empty user when name does not exist", func(t *testing.T) {
		s := NewStorage()
		s.users = append(s.users, User{FirstName: "john", LastName: "doe", Age: 30})

		user, err := s.UserByName("jane")

		assert.NoError(t, err)
		assert.Equal(t, User{}, user)
	})
}
