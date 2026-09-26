package file

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-dyndns/internal/ports"
)

func TestNewFileDNSRepository(t *testing.T) {
	t.Run("Creates missing storage directory", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "nested", "storage")
		filePath := filepath.Join(dir, "dns.json")

		repo, err := NewFileDNSRepository(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, repo)
		assert.DirExists(t, dir)
	})

	t.Run("Storage path unusable", func(t *testing.T) {
		blockingFile := filepath.Join(t.TempDir(), "not-a-dir")
		assert.NoError(t, os.WriteFile(blockingFile, []byte("x"), 0644))

		_, err := NewFileDNSRepository(filepath.Join(blockingFile, "dns.json"))

		assert.Error(t, err)
	})
}

func TestFileDNSRepositoryPing(t *testing.T) {
	ctx := context.Background()

	t.Run("Storage directory reachable", func(t *testing.T) {
		repo, err := NewFileDNSRepository(filepath.Join(t.TempDir(), "dns.json"))
		assert.NoError(t, err)

		assert.NoError(t, repo.Ping(ctx))
	})

	t.Run("Storage directory removed after startup", func(t *testing.T) {
		dir := t.TempDir()
		repo, err := NewFileDNSRepository(filepath.Join(dir, "dns.json"))
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(dir))

		assert.Error(t, repo.Ping(ctx))
	})
}

func TestFileDNSRepositorySaveAndFind(t *testing.T) {
	ctx := context.Background()

	t.Run("Find returns nil for missing domain", func(t *testing.T) {
		repo, err := NewFileDNSRepository(filepath.Join(t.TempDir(), "dns.json"))
		assert.NoError(t, err)

		record, err := repo.Find(ctx, "missing.example.com")

		assert.NoError(t, err)
		assert.Nil(t, record)
	})

	t.Run("Save inserts then Find serves from cache without touching disk", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "dns.json")
		repo, err := NewFileDNSRepository(filePath)
		assert.NoError(t, err)

		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))

		record, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, "home.example.com", record.Domain)
		assert.True(t, net.ParseIP("203.0.113.42").Equal(record.IP))

		assert.FileExists(t, filePath)
	})

	t.Run("Save updates existing domain in place", func(t *testing.T) {
		repo, err := NewFileDNSRepository(filepath.Join(t.TempDir(), "dns.json"))
		assert.NoError(t, err)

		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))
		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("198.51.100.7")}))

		record, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		assert.True(t, net.ParseIP("198.51.100.7").Equal(record.IP))
	})

	t.Run("Find result is a copy safe from external mutation", func(t *testing.T) {
		repo, err := NewFileDNSRepository(filepath.Join(t.TempDir(), "dns.json"))
		assert.NoError(t, err)

		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))

		record, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		record.IP = net.ParseIP("10.0.0.1")

		again, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		assert.True(t, net.ParseIP("203.0.113.42").Equal(again.IP))
	})

	t.Run("Records persisted to disk are reloaded on restart", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "dns.json")
		repo, err := NewFileDNSRepository(filePath)
		assert.NoError(t, err)
		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))

		reopened, err := NewFileDNSRepository(filePath)
		assert.NoError(t, err)

		record, err := reopened.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.True(t, net.ParseIP("203.0.113.42").Equal(record.IP))
	})

	t.Run("Persists one NDJSON record per line, sorted by domain", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "dns.ndjson")
		repo, err := NewFileDNSRepository(filePath)
		assert.NoError(t, err)

		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "office.example.com", IP: net.ParseIP("198.51.100.7")}))
		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))

		raw, err := os.ReadFile(filePath)
		assert.NoError(t, err)

		lines := strings.Split(strings.TrimRight(string(raw), "\n"), "\n")
		assert.Len(t, lines, 2)
		assert.Contains(t, lines[0], `"domain":"home.example.com"`)
		assert.Contains(t, lines[1], `"domain":"office.example.com"`)
	})
}
