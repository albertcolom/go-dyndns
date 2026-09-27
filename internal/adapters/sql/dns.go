package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	"go-dyndns/internal/ports"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Save(ctx context.Context, dns *ports.Dns) error {
	query := `REPLACE INTO dns_records (domain, ip, created_at, updated_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, dns.Domain, dns.IP.String(), dns.CreatedAt.Format(time.RFC3339), dns.UpdatedAt.Format(time.RFC3339))
	return err
}

func (r *SQLRepository) Delete(ctx context.Context, domain string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM dns_records WHERE domain = ?`, domain)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ports.ErrDomainNotFound
	}

	return nil
}

func (r *SQLRepository) Find(ctx context.Context, domain string) (*ports.Dns, error) {
	var ip, createdAt, updatedAt string
	query := `SELECT ip, created_at, updated_at FROM dns_records WHERE domain = ?`
	err := r.db.QueryRowContext(ctx, query, domain).Scan(&ip, &createdAt, &updatedAt)
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

	record := &ports.Dns{Domain: domain, IP: parsedIP}
	record.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	record.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return record, nil
}
