package main

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	env := os.Getenv("ENV")
	db := os.Getenv(fmt.Sprintf("%s_%s", env, "DATABASE_URL"))
	fmt.Println("Connecting to database:", db)

	m, err := migrate.New("file://migrations", db)
	if err != nil {
		panic(err)
	}

	m.Up()
}
