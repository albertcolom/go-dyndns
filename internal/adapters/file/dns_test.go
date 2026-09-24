package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
