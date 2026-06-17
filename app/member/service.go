package member

type service struct {
	storage Storager
	clock   Clock
}

func NewService(storage Storager, clock Clock) *service {
	return &service{
		storage: storage,
		clock:   clock,
	}
}
