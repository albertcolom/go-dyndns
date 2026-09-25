package main

import (
	"context"
	"go-dyndns/config"
	server "go-dyndns/internal/adapters/dns"
	"go-dyndns/internal/adapters/file"
	"go-dyndns/internal/adapters/http"
	"go-dyndns/internal/adapters/http/handler"
	"go-dyndns/internal/adapters/logger"
	"go-dyndns/internal/adapters/sql"
	"go-dyndns/internal/ports"
	dnsservice "go-dyndns/internal/service"
	"log"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	l := logger.NewSlogLogger(cfg.Log.Level)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsnService := dnsservice.NewDSNService()
	dsn, err := dsnService.ParseDSN(cfg.Db.Dsn)
	if err != nil {
		l.Error(ctx, "Failed to parse DSN", "component", "APP", "error", err)
		os.Exit(1)
	}

	var repo ports.DNSRepository
	var healthChecker ports.HealthChecker

	switch dsn.Driver {
	case "file":
		fileRepo, err := file.NewFileDNSRepository(dsn.DataSource)
		if err != nil {
			l.Error(ctx, "Failed to initialize file repository", "component", "APP", "error", err)
			os.Exit(1)
		}
		repo = fileRepo
		healthChecker = fileRepo

	case "sqlite3", "mysql":
		dbClient, err := sql.NewSqlClient(dsn)
		if err != nil {
			l.Error(ctx, "Failed to initialize database client", "component", "APP", "error", err)
			os.Exit(1)
		}
		defer func() {
			if err := dbClient.Close(); err != nil {
				l.Error(ctx, "Database client close error", "component", "APP", "error", err)
			}
		}()

		repo = sql.NewSQLRepository(dbClient.DB)
		healthChecker = dbClient

	default:
		l.Error(ctx, "Unsupported driver", "component", "APP", "driver", dsn.Driver)
		os.Exit(1)
	}

	service := dnsservice.NewDNSService(repo)

	dnsHandler := server.NewDnsHandler(service, l)
	dnsServer := server.NewDnsServer(dnsHandler, cfg.Dns.Addr, cfg.Dns.Net, l)
	dnsErrChan := StartDNSServer(ctx, dnsServer, l)

	httpHandler := handler.NewHandler(service)
	healthHandler := handler.NewHealthHandler(map[string]ports.HealthChecker{
		"Database":   healthChecker,
		"DNS server": dnsServer,
	})
	httpServer := http.NewHTTPServer(httpHandler, healthHandler, cfg.Http.Addr, cfg.Http.Token, l)
	httpErrChan := StartHTTPServer(ctx, httpServer, l)

	WaitForShutdown(ctx, cancel, dnsServer, httpServer, httpErrChan, dnsErrChan, l)
}
