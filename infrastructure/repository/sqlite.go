package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"

	"go-dyndns/domain"
)

type SQLiteDNSRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteDNSRepository {
	return &SQLiteDNSRepository{db: db}
}

func (r SQLiteDNSRepository) Save(record domain.DNSRecord) error {
	_, err := r.db.Exec("REPLACE INTO dns_records (domain, ip) VALUES (?, ?)", record.Domain, record.IP.String())
	return err
}

func (r SQLiteDNSRepository) Find(domainName string) (*domain.DNSRecord, error) {
	var ip string
	err := r.db.QueryRow("SELECT ip FROM dns_records WHERE domain = ?", domainName).Scan(&ip)
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

	return &domain.DNSRecord{Domain: domainName, IP: parsedIP}, nil
}
