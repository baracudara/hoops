package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    var migrationsPath, dbURL string

    flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
    flag.StringVar(&dbURL, "db-url", "", "postgres connection url")
    flag.Parse()

    if migrationsPath == "" {
        panic("migrations-path is required")
    }
    if dbURL == "" {
        panic("db-url is required")
    }

    m, err := migrate.New(
        "file://"+migrationsPath,
        dbURL,
    )
    if err != nil {
        panic(err)
    }

    if err := m.Up(); err != nil {
        if errors.Is(err, migrate.ErrNoChange) {
            fmt.Println("no migrations to apply")
            return
        }
        panic(err)
    }

    fmt.Println("migrations applied")
}