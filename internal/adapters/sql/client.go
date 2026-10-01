package sql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"go-dyndns/internal/ports"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 5 * time.Minute
)

type Client struct {
	DB     *sql.DB
	Driver string
}

func NewSqlClient(dsn *ports.DSN) (*Client, error) {
	db, err := sql.Open(dsn.Driver, dsn.DataSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Client{DB: db, Driver: dsn.Driver}, nil
}

func (c *Client) Close() error {
	return c.DB.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.DB.PingContext(ctx)
}
