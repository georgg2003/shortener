# Shortener

A high-performance URL shortening service built with Go.

Supports multiple storage backends, JWT authentication, gzip compression, and batch operations.

## Features

- Shorten single URLs and process batch requests
- Redirect by short ID with 307 status
- JWT-based authentication (cookie)
- Gzip request/response compression
- Pluggable storage: in-memory, file, PostgreSQL
- Soft delete (user can delete their own URLs)
- pprof profiling endpoint; memory allocations optimized after profiling

## Stack

Go · PostgreSQL · pgx · Chi · JWT · golang-migrate · gRPC · pprof

## Running locally

```bash
# Start PostgreSQL
docker compose up -d postgres

# Run the service
go run ./cmd/shortener \
  -a "localhost:8080" \
  -b "http://localhost:8080" \
  -d "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"
```

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/` | Shorten a URL (plain text) |
| GET | `/{id}` | Redirect to original URL |
| POST | `/api/shorten` | Shorten a URL (JSON) |
| POST | `/api/shorten/batch` | Shorten multiple URLs |
| GET | `/api/user/urls` | List user's URLs |
| DELETE | `/api/user/urls` | Delete user's URLs |
| GET | `/ping` | DB health check |

gRPC API is defined in [api/shortener.proto](api/shortener.proto)

## Configuration

| Flag | Env | Description |
|------|-----|-------------|
| `-a` | `SERVER_ADDRESS` | HTTP server address |
| `-b` | `BASE_URL` | Base URL for short links |
| `-d` | `DATABASE_DSN` | PostgreSQL DSN |
| `-f` | `FILE_STORAGE_PATH` | File storage path (fallback) |
