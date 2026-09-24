package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"go-dyndns/internal/port"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Save(ctx context.Context, dns *port.Dns) error {
	query := `REPLACE INTO dns_records (domain, ip) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, dns.Domain, dns.IP.String())
	return err
}

func (r *SQLRepository) Find(ctx context.Context, domain string) (*port.Dns, error) {
	var ip string
	query := `SELECT ip FROM dns_records WHERE domain = ?`
	err := r.db.QueryRowContext(ctx, query, domain).Scan(&ip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP in database")
	}

	return &port.Dns{Domain: domain, IP: parsedIP}, nil
}
