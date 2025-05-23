package main

import (
	"context"
	"go-dyndns/config"
	server "go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/internal/adapters/repository"
	"go-dyndns/internal/core/dns"
	"go-dyndns/pkg/db"
	"go-dyndns/pkg/logger"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	l, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn, err := db.ParseDSN(cfg.Db.Dsn)
	if err != nil {
		l.Fatal("APP", "Failed to parse DSN", logger.Field{Key: "error", Value: err})
	}

	var repo dns.Repository

	switch dsn.Driver {
	case "file":
		repo = repository.NewFileDNSRepository(dsn.DataSource)

	case "sqlite3", "mysql":
		dbClient, err := db.NewSqlClient(dsn)
		if err != nil {
			l.Fatal("APP", "Failed to initialize database client", logger.Field{Key: "error", Value: err})
		}
		defer func() {
			if err := dbClient.Close(); err != nil {
				l.Error("APP", "Database client close error", logger.Field{Key: "error", Value: err})
			}
		}()

		repo = repository.NewSQLiteDNSRepository(dbClient.DB)

	default:
		l.Fatal("APP", "Unsupported driver", logger.Field{Key: "driver", Value: dsn.Driver})
	}

	service := dns.NewService(repo)

	dnsHandler := server.NewDnsHandler(service, l)
	dnsServer := server.NewDnsServer(dnsHandler, cfg.Dns.Addr, cfg.Dns.Net, l)
	dnsErrChan := StartDNSServer(dnsServer, l)

	httpHandler := http.NewHandler(service)
	httpServer := http.NewHTTPServer(httpHandler, cfg.Http.Addr, cfg.Http.Token, l)
	httpErrChan := StartHTTPServer(httpServer, l)

	WaitForShutdown(cancel, dnsServer, httpServer, httpErrChan, dnsErrChan, l)
}
