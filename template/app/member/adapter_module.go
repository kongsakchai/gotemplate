package member

import "github.com/jmoiron/sqlx"

type External struct {
	DB    *sqlx.DB
	Clock Clock
}

type Module struct {
	Handler *handler
}

func NewModule(adp External) *Module {
	st := NewStorage(adp.DB)
	sv := NewService(st, adp.Clock)
	h := NewHandler(sv)

	return &Module{Handler: h}
}
