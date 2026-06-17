package database

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMySQL(datasource string) (*sqlx.DB, func(context.Context) error) {
	return newDatabase("mysql", datasource)
}
