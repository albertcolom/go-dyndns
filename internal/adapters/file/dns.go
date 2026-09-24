package file

import (
	"context"
	"encoding/json"
	"fmt"
	"go-dyndns/internal/port"
	"os"
	"path/filepath"
	"sync"
)

type FileDNSRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewFileDNSRepository(filePath string) *FileDNSRepository {
	return &FileDNSRepository{filePath: filePath}
}

func (r *FileDNSRepository) Save(ctx context.Context, dns *port.Dns) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.loadRecords()
	if err != nil {
		return err
	}

	updated := false
	for i, record := range data {
		if record.Domain == dns.Domain {
			data[i].IP = dns.IP
			updated = true
			break
		}
	}

	if !updated {
		data = append(data, dns)
	}

	return r.saveRecords(data)
}

func (r *FileDNSRepository) Find(ctx context.Context, domain string) (*port.Dns, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	records, err := r.loadRecords()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.Domain == domain {
			return record, nil
		}
	}

	return nil, nil
}

func (r *FileDNSRepository) loadRecords() ([]*port.Dns, error) {
	var records []*port.Dns
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return records, nil
	}

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	if err := json.Unmarshal(file, &records); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return records, nil
}

func (r *FileDNSRepository) saveRecords(records []*port.Dns) error {
	content, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(r.filePath)
	tmp, err := os.CreateTemp(dir, filepath.Base(r.filePath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	if err := os.Chmod(tmpPath, 0644); err != nil {
		return fmt.Errorf("failed to set permissions on temp file: %w", err)
	}

	if err := os.Rename(tmpPath, r.filePath); err != nil {
		return fmt.Errorf("failed to replace JSON file: %w", err)
	}

	return nil
}
