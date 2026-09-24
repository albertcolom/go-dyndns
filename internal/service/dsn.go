package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"go-dyndns/internal/ports"
)

type dsnService struct{}

func NewDSNService() ports.DSNService {
	return &dsnService{}
}

func (s *dsnService) ParseDSN(raw string) (*ports.DSN, error) {
	parse := strings.Index(raw, "://")
	if parse == -1 {
		return nil, fmt.Errorf("invalid DSN schema: %s", raw)
	}

	driver := raw[:parse]
	dataSource := raw[parse+3:]

	normalizedDriver, ok := ports.SchemeAliases[strings.ToLower(driver)]
	if !ok {
		return nil, fmt.Errorf("unsupported DSN driver: %s", driver)
	}

	normalizedURL := fmt.Sprintf("%s://%s", normalizedDriver, dataSource)

	valid := validate(driver, normalizedURL)
	if !valid {
		return nil, fmt.Errorf("invalid DSN format: %s", raw)
	}

	return &ports.DSN{
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
