package database

import "database/sql"

func newDatabase(driverName string, dataSourceName string) (*sql.DB, func()) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic("Connect to database error: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic("Ping database error: " + err.Error())
	}

	close := func() {
		_ = db.Close()
	}

	return db, close
}
