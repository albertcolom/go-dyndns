[![Go Report Card](https://goreportcard.com/badge/github.com/albertcolom/go-dyndns)](https://goreportcard.com/report/github.com/albertcolom/go-dyndns)
[![Test Status](https://github.com/albertcolom/go-dyndns/actions/workflows/ci.yml/badge.svg)](https://github.com/albertcolom/go-dyndns/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/albertcolom/go-dyndns)](https://github.com/albertcolom/go-dyndns/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/albertcolom/go-dyndns)](https://github.com/albertcolom/go-dyndns/issues)
[![Go Version](https://img.shields.io/badge/go-%3E=1.26-blue)](https://golang.org/doc/go1.26)

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
  token: "change-me"  # API token for authentication

dns:
  addr: ":53"     # DNS server listen address
  net: "udp"      # Network protocol (typically "udp")

db:
  dsn: "sqlite3://./app.db"  # Storage backend (see Supported Database section)

log:
  level: "info"   # Log level: debug, info, warn, error
```
Every configuration value in the YAML can be overridden by setting an environment variable.

| YAML Key     | Environment Variable | Description                          |
|--------------|----------------------|---------------------------------------|
| `http.addr`  | `HTTP_ADDR`          | HTTP server listen address            |
| `http.token` | `HTTP_TOKEN`         | API authentication token              |
| `dns.addr`   | `DNS_ADDR`           | DNS server listen address             |
| `dns.net`    | `DNS_NET`            | DNS protocol (e.g., `udp`, `tcp`)     |
| `db.dsn`     | `DB_DSN`             | DSN connection string                 |
| `log.level`  | `LOG_LEVEL`          | Log level: `debug`, `info`, `warn`, `error` |


## 🗄️ Supported Database
The go-dyndns service supports multiple database backends through a unified DSN (Data Source Name) format. You can configure your backend via the `db.dsn` field in your YAML config file or override it using the `DB_DSN` environment variable.

Below is a summary of supported drivers:

| Driver    | DSN Format Example                                                        | Description                  | Go Driver Package                         | Multi-replica safe |
|-----------|----------------------------------------------------------------------------|------------------------------|--------------------------------------------|:-------------------:|
| `file`    | `file://./app.ndjson`                                                     | NDJSON file storage on disk (one record per line, in-memory cache) | _Built-in (no external dependency)_       | ❌ |
| `sqlite`  | `sqlite://./app.db`                                                       | Alias of `sqlite3`             |  | ❌ |
| `sqlite3` | `sqlite3://./app.db`                                                      | Lightweight SQLite database  | [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) | ❌ |
| `mysql`   | `mysql://root:root@tcp(localhost:3306)/app?tls=false`         | MySQL or MariaDB SQL backend | [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql) | ✅ |

> **Running more than one instance (e.g. multiple Kubernetes pods)?** Use the `mysql` backend. `file` and `sqlite`/`sqlite3` store data on the local disk of a single instance, guarded only by an in-process lock — with several replicas, each one would read/write its own separate copy of the data (or race on a shared volume), so DNS records would silently go out of sync between instances. Scale those backends by keeping a single replica.

## 🔄 Router / Firmware DDNS Clients (FritzBox, ddclient, inadyn)
Besides the JSON `/v1/domains` API, go-dyndns exposes a `/nic/update` endpoint
implementing the classic **dyndns2** protocol used by built-in DDNS clients in
routers and firmware — so you can point them at go-dyndns directly, with no
custom scripting.

- **Method/path:** `GET /nic/update?hostname=<domain>&myip=<ip>`
- **Auth:** HTTP Basic Auth — any username, with the password set to your `http.token`
- **Response:** plain text — `good <ip>` (record created/changed), `nochg <ip>` (already up to date), or an error word (`notfqdn`, `badauth`, `911`)

Unlike `/v1/domains`, this endpoint is an upsert: the first update for a
hostname creates it, later ones update it — there's no separate create step.

### AVM FritzBox setup
In FRITZ!Box: **Internet → Permit Access → DynDNS**, choose **User-defined**, and set:

| Field | Value |
|---|---|
| Update-URL | `http://<your-server>/nic/update?hostname=<domain>&myip=<ipaddr>` |
| Domain name | your domain, e.g. `home.example.com` |
| Username | anything (ignored) |
| Password | your `http.token` |

> IPv6 isn't supported yet (only `<ipaddr>`/IPv4); `<ip6addr>` placeholders are not handled.

## 🧬 Database Migrations
This app supports database migrations (e.g., for SQLite/MySQL/MariaDB) using a migration tool in `./cmd/migrations`

Available commands:
```bash
make migrate-up         # Apply all up migrations
make migrate-down       # Roll back the last migration
make migrate-version    # Show current migration version
```
