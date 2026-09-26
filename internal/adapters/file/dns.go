package file

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"go-dyndns/internal/ports"
)

type FileDNSRepository struct {
	filePath string
	mu       sync.RWMutex
	records  map[string]*ports.Dns
}

func NewFileDNSRepository(filePath string) (*FileDNSRepository, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	r := &FileDNSRepository{filePath: filePath}

	loaded, err := r.loadRecords()
	if err != nil {
		return nil, err
	}

	r.records = make(map[string]*ports.Dns, len(loaded))
	for _, record := range loaded {
		r.records[record.Domain] = record
	}

	return r, nil
}

func (r *FileDNSRepository) Save(_ context.Context, dns *ports.Dns) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record := cloneDns(dns)

	records := make([]*ports.Dns, 0, len(r.records)+1)
	for domain, existing := range r.records {
		if domain == dns.Domain {
			continue
		}
		records = append(records, existing)
	}
	records = append(records, record)

	sort.Slice(records, func(i, j int) bool { return records[i].Domain < records[j].Domain })

	if err := r.saveRecords(records); err != nil {
		return err
	}

	r.records[dns.Domain] = record

	return nil
}

func (r *FileDNSRepository) Find(_ context.Context, domain string) (*ports.Dns, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.records[domain]
	if !ok {
		return nil, nil
	}

	return cloneDns(record), nil
}

func cloneDns(dns *ports.Dns) *ports.Dns {
	return &ports.Dns{Domain: dns.Domain, IP: append(net.IP(nil), dns.IP...)}
}

func (r *FileDNSRepository) loadRecords() ([]*ports.Dns, error) {
	file, err := os.Open(r.filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read NDJSON file: %w", err)
	}
	defer file.Close()

	var records []*ports.Dns
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var record ports.Dns
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, fmt.Errorf("failed to parse NDJSON line: %w", err)
		}
		records = append(records, &record)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read NDJSON file: %w", err)
	}

	return records, nil
}

func (r *FileDNSRepository) saveRecords(records []*ports.Dns) error {
	var buf bytes.Buffer
	for _, record := range records {
		line, err := json.Marshal(record)
		if err != nil {
			return err
		}
		buf.Write(line)
		buf.WriteByte('\n')
	}
	content := buf.Bytes()

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
		return fmt.Errorf("failed to replace NDJSON file: %w", err)
	}

	return nil
}

func (r *FileDNSRepository) Ping(_ context.Context) error {
	dir := filepath.Dir(r.filePath)
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("storage directory unreachable: %w", err)
	}
	return nil
}
