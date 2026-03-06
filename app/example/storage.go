package example

type storage struct {
	users []User
}

func NewStorage() *storage {
	users := make([]User, 0)
	return &storage{users: users}
}

func (s *storage) CreateUser(user User) error {
	s.users = append(s.users, user)
	return nil
}

func (s *storage) UserByName(name string) (User, error) {
	for _, user := range s.users {
		if user.FirstName == name {
			return user, nil
		}
	}
	return User{}, nil
}
