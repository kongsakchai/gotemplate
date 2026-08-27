package postgresdb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/database/sqldb"
	_ "github.com/lib/pq"
)

func New(datasource string) (*sqlx.DB, func(context.Context) error) {
	return sqldb.New("postgres", datasource)
}
