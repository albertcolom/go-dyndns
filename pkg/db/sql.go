package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	DB *sql.DB
}

func NewSqlClient(dsn string) (*Client, error) {
	driver, source, found := strings.Cut(dsn, "://")
	if !found {
		return nil, fmt.Errorf("invalid DSN format: missing scheme (://): %s", dsn)
	}

	switch driver {
	case "sqlite", "sqlite3":
		driver = "sqlite3"
	case "postgres", "postgresql":
		driver = "postgres"
	case "mysql":
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Client{DB: db}, nil
}

func (c *Client) Close() error {
	return c.DB.Close()
}
