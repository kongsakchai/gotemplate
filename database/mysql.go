package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kongsakchai/gotemplate/config"
)

var mysqlDB *sql.DB

func NewMySQL(cfg config.Database) (*sql.DB, func()) {
	db, err := sql.Open("mysql", cfg.URL)

	if err != nil {
		panic("Connect to database error: " + err.Error())
	}

	close := func() {
		_ = db.Close()
	}

	mysqlDB = db
	return db, close
}

func IsMySQLReady() bool {
	return mysqlDB.Ping() == nil
}
