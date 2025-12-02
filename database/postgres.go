package database

import (
	"database/sql"

	"github.com/kongsakchai/gotemplate/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.Database) (*sql.DB, func() error) {
	return newDatabase("postgres", cfg.URL)
}
