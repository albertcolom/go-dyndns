package db

import (
	"fmt"
	"net/url"
	"strings"
)

var schemeAliases = map[string]string{
	"sqlite":     "sqlite3",
	"sqlite3":    "sqlite3",
	"postgres":   "postgres",
	"postgresql": "postgres",
	"mysql":      "mysql",
}

type DSN struct {
	Driver     string
	DataSource string
	Raw        string
	Normalized string
}

func ParseDSN(raw string) (*DSN, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" {
		return nil, fmt.Errorf("invalid DSN format: %s", raw)
	}

	normalizedDriver, ok := schemeAliases[strings.ToLower(u.Scheme)]
	if !ok {
		return nil, fmt.Errorf("unsupported DSN driver: %s", u.Scheme)
	}

	normalizedURL := *u
	normalizedURL.Scheme = normalizedDriver
	dataSource := strings.TrimPrefix(raw, u.Scheme+"://")

	return &DSN{
		Driver:     normalizedDriver,
		DataSource: dataSource,
		Raw:        raw,
		Normalized: normalizedURL.String(),
	}, nil
}
