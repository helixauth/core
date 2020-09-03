package main

// func main() {
// 	panic("Not implemented")
// }

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/helixauth/helix/cfg"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	// Parse arguments
	if len(os.Args) != 2 {
		panic(errors.New("Syntax: go run bin/migrate/main.go [up/down]"))
	}
	direction := strings.ToLower(os.Args[1])
	if direction != "up" && direction != "down" {
		panic(errors.New("Syntax: go run bin/migrate/main.go [up/down]"))
	}

	// Get the DB's admin password
	fmt.Print("Sysadmin password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err.Error())
	}
	log.Println()

	// Connect to the database
	cfg := config.New()
	connInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", "sysadmin", string(password), cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	m, err := migrate.New("file://bin/migrate/sql", connInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate
	if direction == "up" {
		err = m.Up()
	} else if direction == "down" {
		err = m.Down()
	}
	if err != nil {
		log.Fatal(err)
	}

	version, _, _ := m.Version()
	fmt.Printf("Migrated database %v to version %v\n\n", direction, version)
}
