package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
		return nil, err
	}

	return &config, nil
}
