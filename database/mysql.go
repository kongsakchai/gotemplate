package database

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kongsakchai/gotemplate/config"
)

func NewMySQL(cfg config.Database) (*sql.DB, func(context.Context) error) {
	return newDatabase("mysql", cfg.URL)
}
