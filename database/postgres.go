package database

import (
	"database/sql"

	"github.com/kongsakchai/gotemplate/config"
	_ "github.com/lib/pq"
)

var postgresDB *sql.DB

func NewPostgres(cfg config.Database) (*sql.DB, func()) {
	db, err := sql.Open("postgres", cfg.URL)

	if err != nil {
		panic("Connect to database error: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		panic("Ping database error: " + err.Error())
	}

	close := func() {
		_ = db.Close()
	}

	postgresDB = db
	return db, close
}

func IsPostgresReady() bool {
	return postgresDB.Ping() == nil
}
