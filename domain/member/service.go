package member

import "time"

//mockery:generate: true
type Clock interface {
	Now() time.Time
}

type ServiceDeps struct {
	Storage Storager
	Clock   Clock
}

type service struct {
	storage Storager
	clock   Clock
}

func NewService(deps ServiceDeps) *service {
	return &service{
		storage: deps.Storage,
		clock:   deps.Clock,
	}
}
