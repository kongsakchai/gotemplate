package database

import (
	"database/sql"
)

func newDatabase(driverName string, dataSourceName string) (*sql.DB, func() error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic("Connect to database error: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic("Ping database error: " + err.Error())
	}

	close := func() error {
		return db.Close()
	}

	return db, close
}
