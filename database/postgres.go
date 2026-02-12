package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.Database) (*sqlx.DB, func(context.Context) error) {
	return newDatabase("postgres", cfg.URL)
}
