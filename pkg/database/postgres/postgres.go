package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/database"
	_ "github.com/lib/pq"
)

func New(datasource string) (*sqlx.DB, func(context.Context) error) {
	return database.NewDatabase("postgres", datasource)
}
