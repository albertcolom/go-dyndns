package sql

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-dyndns/internal/ports"
)

func TestNewSqlClient(t *testing.T) {
	t.Run("connects to sqlite3 and applies pool settings", func(t *testing.T) {
		dbPath := filepath.Join(t.TempDir(), "app.db")

		client, err := NewSqlClient(&ports.DSN{Driver: "sqlite3", DataSource: dbPath})
		assert.NoError(t, err)
		assert.NotNil(t, client)
		t.Cleanup(func() { client.Close() })

		stats := client.DB.Stats()
		assert.Equal(t, maxOpenConns, stats.MaxOpenConnections)
	})

	t.Run("fails for an unsupported driver", func(t *testing.T) {
		_, err := NewSqlClient(&ports.DSN{Driver: "unsupported", DataSource: "x"})
		assert.Error(t, err)
	})
}
