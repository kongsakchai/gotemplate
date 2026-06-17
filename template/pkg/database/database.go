package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func newDatabase(driverName string, dataSourceName string) (*sqlx.DB, func(context.Context) error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		panic("Connect to database error: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic("Ping database error: " + err.Error())
	}

	close := func(_ context.Context) error {
		return db.Close()
	}

	return db, close
}
