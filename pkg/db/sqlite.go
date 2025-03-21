package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	DB *sql.DB
}

func NewSQLiteClient(dbPath string) (*Client, error) {
	db, err := sql.Open("sqlite3", dbPath)
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
