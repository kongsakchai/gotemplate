package sqlitedb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/database/sqldb"
	_ "modernc.org/sqlite"
)

func New(datasource string) (*sqlx.DB, func(context.Context) error) {
	return sqldb.New("sqlite", datasource)
}
