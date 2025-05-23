package main

import (
	"errors"
	"flag"
	"fmt"
	"go-dyndns/config"
	"go-dyndns/pkg/db"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn, err := db.ParseDSN(cfg.Db.Dsn)
	if err != nil {
		log.Fatalf("Failed to parse dsn: %v", err)
	}

	m, err := migrate.New("file://database/migrations", dsn.Normalized)
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v", err)
	}

	upFlag := flag.Bool("up", false, "Apply all available migrations")
	downFlag := flag.Bool("down", false, "Rollback the most recent migration")
	versionFlag := flag.Bool("version", false, "Show the current migration version")
	flag.Parse()

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
