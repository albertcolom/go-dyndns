package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDSNSuccess(t *testing.T) {
	tests := []struct {
		raw                string
		expectedDriver     string
		expectedDataSource string
		expectedNormalized string
	}{
		{
			raw:                "sqlite://file.db",
			expectedDriver:     "sqlite3",
			expectedDataSource: "file.db",
			expectedNormalized: "sqlite3://file.db",
		},
		{
			raw:                "postgres://user:password@localhost:5432/dbname",
			expectedDriver:     "postgres",
			expectedDataSource: "user:password@localhost:5432/dbname",
			expectedNormalized: "postgres://user:password@localhost:5432/dbname",
		},
		{
			raw:                "mysql://user:password@localhost:3306/dbname",
			expectedDriver:     "mysql",
			expectedDataSource: "user:password@localhost:3306/dbname",
			expectedNormalized: "mysql://user:password@localhost:3306/dbname",
		},
		{
			raw:                "postgresql://user@localhost:5432/dbname",
			expectedDriver:     "postgres",
			expectedDataSource: "user@localhost:5432/dbname",
			expectedNormalized: "postgres://user@localhost:5432/dbname",
		},
	}

	for _, tt := range tests {
		t.Run("Parsing "+tt.raw, func(t *testing.T) {
			dsn, err := ParseDSN(tt.raw)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedDriver, dsn.Driver)
			assert.Equal(t, tt.expectedDataSource, dsn.DataSource)
			assert.Equal(t, tt.expectedNormalized, dsn.Normalized)
		})
	}
}

func TestParseDSNError(t *testing.T) {
	tests := []struct {
		raw           string
		expectedError string
	}{
		{
			raw:           "unsupported://user:password@localhost:3306/dbname",
			expectedError: "unsupported DSN driver: unsupported",
		},
		{
			raw:           "invalidurl",
			expectedError: "invalid DSN format: invalidurl",
		},
		{
			raw:           "ftp://user:password@localhost:3306/dbname",
			expectedError: "unsupported DSN driver: ftp",
		},
	}

	for _, tt := range tests {
		t.Run("Parsing "+tt.raw, func(t *testing.T) {
			_, err := ParseDSN(tt.raw)
			assert.Error(t, err)
			assert.Equal(t, err.Error(), tt.expectedError)
		})
	}
}
