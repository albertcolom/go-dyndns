package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"go-dyndns/internal/core"
	"os"
	"sync"
)

type FileDNSRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewFileDNSRepository(filePath string) *FileDNSRepository {
	return &FileDNSRepository{filePath: filePath}
}

func (r *FileDNSRepository) Save(ctx context.Context, dns *core.Dns) error {
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

func (r *FileDNSRepository) Find(ctx context.Context, domain string) (*core.Dns, error) {
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

func (r *FileDNSRepository) loadRecords() ([]*core.Dns, error) {
	var records []*core.Dns
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

func (r *FileDNSRepository) saveRecords(records []*core.Dns) error {
	content, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, content, 0644)
}
