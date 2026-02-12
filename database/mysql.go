package database

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/config"
)

func NewMySQL(cfg config.Database) (*sqlx.DB, func(context.Context) error) {
	return newDatabase("mysql", cfg.URL)
}
