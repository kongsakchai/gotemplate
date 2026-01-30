package database

import (
	"context"
	"database/sql"

	"github.com/kongsakchai/gotemplate/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.Database) (*sql.DB, func(context.Context) error) {
	return newDatabase("postgres", cfg.URL)
}
