package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"go-dyndns/internal/core/dns"
)

type SQLiteDNSRepository struct {
	db *sql.DB
}

func NewSQLiteDNSRepository(db *sql.DB) *SQLiteDNSRepository {
	return &SQLiteDNSRepository{db: db}
}

func (r *SQLiteDNSRepository) Save(ctx context.Context, dns *dns.Dns) error {
	query := `REPLACE INTO dns_records (domain, ip) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, dns.Domain, dns.IP.String())
	return err
}

func (r *SQLiteDNSRepository) Find(ctx context.Context, domain string) (*dns.Dns, error) {
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

	return &dns.Dns{Domain: domain, IP: parsedIP}, nil
}
