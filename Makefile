# Имя бинаря
BINARY := shortener

# Путь к миграциям
MIGRATIONS_DIR := ./migrations

# Параметры базы
DB_URL := "postgres://shortener:password@localhost:5432/shortener?sslmode=disable"

# Go параметры
GO := go
GOFLAGS :=
GOMOD := $(shell go env GOMOD)
GOPATH := $(shell go env GOPATH)
PKG := ./...

export PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all tidy deps build run test mock migrate-up migrate-down clean

all: build

## ---------------------------
## Dependencies
## ---------------------------

tidy:
	$(GO) mod tidy

deps:
	$(GO) install go.uber.org/mock/mockgen@latest
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## ---------------------------
## Build & Run
## ---------------------------

build:
	$(GO) build -o bin/$(BINARY) $(PKG)

run:
	$(GO) run $(PKG) -d $(DB_URL)

clean:
	rm -rf bin

## ---------------------------
## Testing
## ---------------------------

test:
	$(GO) test $(PKG) -cover

## ---------------------------
## Generate mocks and stuff
## ---------------------------

generate:
	$(GO) generate $(PKG)

## ---------------------------
## Migrations
## ---------------------------

migrate-up:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) up

migrate-down:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) down