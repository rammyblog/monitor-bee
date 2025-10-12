package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/rammyblog/monitor-bee/internal/config"
)

func main() {
	var (
		flags   = flag.NewFlagSet("migrate", flag.ExitOnError)
		dir     = flags.String("dir", "migrations", "directory with migration files")
		verbose = flags.Bool("v", false, "enable verbose mode")
	)

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	if *verbose {
		goose.SetVerbose(true)
	}

	args := flags.Args()
	if len(args) == 0 {
		flags.Usage()
		return
	}

	command := args[0]

	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	switch command {
	case "up":
		if err := goose.Up(db, *dir); err != nil {
			log.Fatalf("failed to run migrations up: %v", err)
		}
	case "down":
		if err := goose.Down(db, *dir); err != nil {
			log.Fatalf("failed to run migrations down: %v", err)
		}
	case "status":
		if err := goose.Status(db, *dir); err != nil {
			log.Fatalf("failed to get migration status: %v", err)
		}
	case "version":
		version, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatalf("failed to get database version: %v", err)
		}
		fmt.Printf("Current version: %d\n", version)
	case "create":
		if len(args) < 2 {
			log.Fatal("create command requires a migration name")
		}
		name := args[1]
		if err := goose.Create(db, *dir, name, "sql"); err != nil {
			log.Fatalf("failed to create migration: %v", err)
		}
	case "reset":
		if err := goose.Reset(db, *dir); err != nil {
			log.Fatalf("failed to reset database: %v", err)
		}
	default:
		log.Fatalf("unknown command: %s", command)
	}
}

func usage() {
	fmt.Print(`Usage: go run cmd/migrate/main.go [OPTIONS] COMMAND

Commands:
    up                   Migrate the DB to the most recent version available
    down                 Roll back the version by 1
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME          Create the scaffolding for a new migration
    reset                Roll back all migrations

Options:
    -dir string          directory with migration files (default "migrations")
    -v                   enable verbose mode

Examples:
    go run cmd/migrate/main.go up
    go run cmd/migrate/main.go down
    go run cmd/migrate/main.go status
    go run cmd/migrate/main.go create add_user_roles
`)
}
