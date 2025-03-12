package repository

import (
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

func (r *SQLiteDNSRepository) Save(dns *dns.Dns) error {
	_, err := r.db.Exec("REPLACE INTO dns_records (domain, ip) VALUES (?, ?)", dns.Domain, dns.IP.String())
	return err
}

func (r *SQLiteDNSRepository) Find(domain string) (*dns.Dns, error) {
	var ip string
	err := r.db.QueryRow("SELECT ip FROM dns_records WHERE domain = ?", domain).Scan(&ip)
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
