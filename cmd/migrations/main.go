package main

import (
	"errors"
	"flag"
	"fmt"
	"go-dyndns/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	upFlag := flag.Bool("up", false, "Apply all available migrations")
	downFlag := flag.Bool("down", false, "Rollback the most recent migration")
	versionFlag := flag.Bool("version", false, "Show the current migration version")
	flag.Parse()

	dbURL := "sqlite3://" + cfg.Sqlite.Path
	migrationsPath := "file://database/migrations"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v", err)
	}

	if *upFlag {
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No new migrations to apply.")
			} else {
				log.Fatalf("migration failed: %v", err)
			}
		} else {
			fmt.Println("Migrations applied successfully.")
		}
		return
	}

	if *downFlag {
		if err := m.Steps(-1); err != nil {
			log.Fatalf("failed to rollback migration: %v", err)
		} else {
			fmt.Println("Most recent migration rolled back successfully.")
		}
		return
	}

	if *versionFlag {
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("failed to get version: %v", err)
		}
		if dirty {
			fmt.Printf("Current version: %d (dirty)\n", version)
		} else {
			fmt.Printf("Current version: %d\n", version)
		}
		return
	}

	fmt.Println("Usage:")
	flag.PrintDefaults()
}
