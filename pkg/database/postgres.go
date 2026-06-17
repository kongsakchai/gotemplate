package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(datasource string) (*sqlx.DB, func(context.Context) error) {
	return newDatabase("postgres", datasource)
}
