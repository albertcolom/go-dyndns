[![Go Report Card](https://goreportcard.com/badge/github.com/albertcolom/go-dyndns)](https://goreportcard.com/report/github.com/albertcolom/go-dyndns)
[![Test Status](https://github.com/albertcolom/go-dyndns/actions/workflows/ci.yml/badge.svg)](https://github.com/albertcolom/go-dyndns/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/albertcolom/go-dyndns)](https://github.com/albertcolom/go-dyndns/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/albertcolom/go-dyndns)](https://github.com/albertcolom/go-dyndns/issues)
[![Go Version](https://img.shields.io/badge/go-%3E=1.23-blue)](https://golang.org/doc/go1.23)

# go-dyndns

## 🧭 go-dyndns
`go-dyndns` is a lightweight and extensible dynamic DNS server written in Go. It provides a simple and efficient way to manage DNS records dynamically through an HTTP API, while serving DNS responses over the standard UDP protocol.

This project is ideal for small home labs, internal networks, IoT devices, or self-hosted services that need to update their DNS entries dynamically — without relying on external providers.

---

## 📦 Configuration
The application is configured using a simple YAML file. You can also override all configuration values using environment variables, which is especially useful in containerized environments like Docker or Kubernetes.

Example `config/config.yaml`
```yaml
http:
  addr: ":8080"   # HTTP server listen address
  token: "a38b721f-e8c8-4cdf-95b0-baa13ee5ddd5"  # API token for authentication

dns:
  addr: ":53"     # DNS server listen address
  net: "udp"      # Network protocol (typically "udp")

db:
  dsn: "sqlite3://./app.db"  # Storage backend (see Supported Database section)
```
Every configuration value in the YAML can be overridden by setting an environment variable.

| YAML Key     | Environment Variable | Description                      |
|--------------|----------------------|----------------------------------|
| `http.addr`  | `HTTP_ADDR`          | HTTP server listen address       |
| `http.token` | `HTTP_TOKEN`         | API authentication token         |
| `dns.addr`   | `DNS_ADDR`           | DNS server listen address        |
| `dns.net`    | `DNS_NET`            | DNS protocol (e.g., `udp`, `tcp`)|
| `db.dsn`     | `DB_DSN`             | DSN connection string            |


## 🗄️ Supported Database
The go-dyndns service supports multiple database backends through a unified DSN (Data Source Name) format. You can configure your backend via the `db.dsn` field in your YAML config file or override it using the `DB_DSN` environment variable.

Below is a summary of supported drivers:

| Driver    | DSN Format Example                                                        | Description                  | Go Driver Package                         |
|-----------|----------------------------------------------------------------------------|------------------------------|-------------------------------------------|
| `file`    | `file://./app.json`                                                       | JSON file storage on disk    | _Built-in (no external dependency)_       |
| `sqlite`  | `sqlite://./app.db`                                                       | Alias of `sqlite3`             |  |
| `sqlite3` | `sqlite3://./app.db`                                                      | Lightweight SQLite database  | [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) |
| `mysql`   | `mysql://root:root@tcp(localhost:3306)/app?tls=false`         | MySQL or MariaDB SQL backend | [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql) |

## 🧬 Database Migrations
This app supports database migrations (e.g., for SQLite/MySQL/MariaDB) using a migration tool in `./cmd/migrations`

Available commands:
```bash
make migrate-up         # Apply all up migrations
make migrate-down       # Roll back the last migration
make migrate-version    # Show current migration version
```

