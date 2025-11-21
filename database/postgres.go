package database

import (
	"database/sql"

	"github.com/kongsakchai/gotemplate/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.Database) (*sql.DB, func()) {
	return newDatabase("postgres", cfg.URL)
}
