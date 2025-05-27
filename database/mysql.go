package database

import (
	"database/sql"
	"log"

	"github.com/kongsakchai/gotemplate/config"
)

func NewMySQL(cfg config.Database) (*sql.DB, func()) {
	db, err := sql.Open("mysql", cfg.URL)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping to database error", err)
	}

	close := func() {
		_ = db.Close()
	}

	return db, close
}
