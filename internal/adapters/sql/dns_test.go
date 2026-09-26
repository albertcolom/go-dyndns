package sql

import (
	"context"
	"database/sql"
	"net"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"go-dyndns/internal/ports"
)

func newTestSQLRepo(t *testing.T) *SQLRepository {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE dns_records(
		domain     VARCHAR(255) PRIMARY KEY,
		ip         VARCHAR(45) NOT NULL,
		created_at VARCHAR(32) NOT NULL DEFAULT '',
		updated_at VARCHAR(32) NOT NULL DEFAULT ''
	)`)
	assert.NoError(t, err)

	return NewSQLRepository(db)
}

func TestSQLRepositorySaveAndFind(t *testing.T) {
	ctx := context.Background()

	t.Run("Find returns nil for missing domain", func(t *testing.T) {
		repo := newTestSQLRepo(t)

		record, err := repo.Find(ctx, "missing.example.com")

		assert.NoError(t, err)
		assert.Nil(t, record)
	})

	t.Run("Save then Find round-trips domain and IP", func(t *testing.T) {
		repo := newTestSQLRepo(t)

		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))

		record, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, "home.example.com", record.Domain)
		assert.True(t, net.ParseIP("203.0.113.42").Equal(record.IP))
	})

	t.Run("Save sets CreatedAt on insert and preserves it across updates while refreshing UpdatedAt", func(t *testing.T) {
		repo := newTestSQLRepo(t)

		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))
		first, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		assert.False(t, first.CreatedAt.IsZero())
		assert.False(t, first.UpdatedAt.IsZero())

		time.Sleep(time.Second)
		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("198.51.100.7")}))
		second, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)

		assert.True(t, first.CreatedAt.Equal(second.CreatedAt))
		assert.True(t, second.UpdatedAt.After(first.UpdatedAt))
	})

	t.Run("Delete returns ErrDomainNotFound for missing domain", func(t *testing.T) {
		repo := newTestSQLRepo(t)

		err := repo.Delete(ctx, "missing.example.com")

		assert.ErrorIs(t, err, ports.ErrDomainNotFound)
	})

	t.Run("Delete removes an existing record", func(t *testing.T) {
		repo := newTestSQLRepo(t)

		assert.NoError(t, repo.Save(ctx, &ports.Dns{Domain: "home.example.com", IP: net.ParseIP("203.0.113.42")}))
		assert.NoError(t, repo.Delete(ctx, "home.example.com"))

		record, err := repo.Find(ctx, "home.example.com")
		assert.NoError(t, err)
		assert.Nil(t, record)
	})
}
