#  CeroDev

**CeroDev** is a modern full-stack platform built with a blazing-fast [Go](https://golang.org) backend and a sleek, minimal frontend powered by [Bun](https://bun.sh) + React.



## 🚀 Features

- 🔐 Auth-protected dashboard
- 🐳 Docker container management
- 🧠 Context-aware logging with trace IDs
- ⚡ Ultra-fast builds via Bun
- 🦫 Go + SQLite + SQLC backend
- 📦 Migrations with `golang-migrate`


##  Prerequisites

- Go 1.24+
- Bun (>= v1.0)
- SQLite
- Docker (for container management)
- `golang-migrate` (for DB migrations)


### Configuration

- Loads from environamen variables and .env


## Database Migrations

Migrations are written in SQL, located in /migrations.

```
migrate -database sqlite3://dev.db -path ./migrations up
```