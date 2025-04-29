package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kongsakchai/gotemplate/config"
)

func NewMySQL(conf config.Config) (*sql.DB, func()) {
	db, err := sql.Open("mysql", conf.Database.MySQLURI)
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
