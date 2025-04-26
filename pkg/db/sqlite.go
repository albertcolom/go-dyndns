package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	DB *sql.DB
}

func NewSQLitClient(dsn string) (*Client, error) {
	driver, source, found := strings.Cut(dsn, "://")
	if !found {
		return nil, fmt.Errorf("invalid DSN format: missing scheme (://): %s", dsn)
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
