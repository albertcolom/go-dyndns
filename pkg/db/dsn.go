package db

import (
	"fmt"
	"net/url"
	"regexp"
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
	parse := strings.Index(raw, "://")
	if parse == -1 {
		return nil, fmt.Errorf("invalid DSN schema: %s", raw)
	}

	driver := raw[:parse]
	dataSource := raw[parse+3:]

	normalizedDriver, ok := schemeAliases[strings.ToLower(driver)]
	if !ok {
		return nil, fmt.Errorf("unsupported DSN driver: %s", driver)
	}

	normalizedURL := fmt.Sprintf("%s://%s", normalizedDriver, dataSource)

	valid := validate(driver, normalizedURL)
	if !valid {
		return nil, fmt.Errorf("invalid DSN format: %s", raw)
	}

	return &DSN{
		Driver:     normalizedDriver,
		DataSource: dataSource,
		Raw:        raw,
		Normalized: normalizedURL,
	}, nil
}

func validate(driver, normalizedURL string) bool {
	if driver == "mysql" {
		regex := `^mysql:\/\/([^:]+):([^@]+)@tcp\(([^:]+):(\d+)\)\/([^?]+)$`
		re := regexp.MustCompile(regex)

		return re.MatchString(normalizedURL)
	}

	u, err := url.Parse(normalizedURL)
	if err != nil || u.Scheme == "" {
		return false
	}

	return true
}
