package mysqldb

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/database/sqldb"
)

func New(datasource string) (*sqlx.DB, func(context.Context) error) {
	return sqldb.New("mysql", datasource)
}
