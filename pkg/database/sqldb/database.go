package sqldb

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

func New(driverName string, dataSourceName string) (*sqlx.DB, func(context.Context) error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		panic("Connect to database error: " + err.Error())
	}

	ctxPing, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctxPing); err != nil {
		panic("Ping database error: " + err.Error())
	}

	close := func(_ context.Context) error {
		return db.Close()
	}

	return db, close
}
