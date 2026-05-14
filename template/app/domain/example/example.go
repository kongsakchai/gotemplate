package example

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type domain struct {
	Handler *handler
}

func New() *domain {
	st := NewStorage()
	h := NewHandler(st)
	return &domain{Handler: h}
}
