package example

import "github.com/kongsakchai/gotemplate/app"

type storage struct {
	users map[string]User
}

func NewStorage() *storage {
	users := make(map[string]User)
	users["john"] = User{Name: "John", Age: 30}

	return &storage{users: users}
}

func (s *storage) UserByName(name string) (User, app.Error) {
	user, exists := s.users[name]
	if !exists {
		return User{}, app.NotFound("4001", "user not found", nil)
	}
	return user, app.Error{}
}

func (s *storage) Users() ([]User, app.Error) {
	var users []User
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, app.Error{}
}
