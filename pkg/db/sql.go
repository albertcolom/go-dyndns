package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	DB     *sql.DB
	Driver string
}

func NewSqlClient(dsn *DSN) (*Client, error) {
	db, err := sql.Open(dsn.Driver, dsn.DataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Client{DB: db, Driver: dsn.Driver}, nil
}

func (c *Client) Close() error {
	return c.DB.Close()
}
