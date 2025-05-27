package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/kongsakchai/gotemplate/config"
)

func NewMySQL(conf config.Database) (*sql.DB, func()) {
	cfg := mysql.Config{
		User:   conf.Username,
		Passwd: conf.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		DBName: conf.DBName,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

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
