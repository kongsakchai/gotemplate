package database

import (
	"database/sql"
	"log"

	"github.com/kongsakchai/gotemplate/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.Database) (*sql.DB, func()) {
	db, err := sql.Open("postgres", cfg.URL)

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
