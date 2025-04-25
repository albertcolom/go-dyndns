package config

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var configFS embed.FS

type Config struct {
	Http   HttpConfig
	Dns    DnsConfig
	Sqlite SqliteConfig
}

type HttpConfig struct {
	Addr  string
	Token string
}

type DnsConfig struct {
	Addr string
	Net  string
}

type SqliteConfig struct {
	Path string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	fileData, err := configFS.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded config: %w", err)
	}

	if err := viper.ReadConfig(bytes.NewReader(fileData)); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
		return nil, err
	}

	return &config, nil
}
