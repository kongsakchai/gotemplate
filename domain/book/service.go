package book

type service struct {
	storage Storager
}

func NewService(storage Storager) *service {
	return &service{
		storage: storage,
	}
}
