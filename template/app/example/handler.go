package example

//mockery:generate: true
type Storager interface {
	Users() ([]User, error)
	UserByName(name string) (User, error)
	CreateUser(user User) error
}

type handler struct {
	storage Storager
}

func NewHandler(storage Storager) *handler {
	return &handler{storage: storage}
}
