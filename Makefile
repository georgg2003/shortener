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
SHORTENER_PKG = ./cmd/shortener/main.go
PKG := ./cmd/shortener
BUILD_VER = v0.1
BUILD_DATE := $(shell date +"%Y/%m/%d %H:%M:%S")
BUILD_COMMIT := $(shell git rev-parse --short HEAD)

export PATH := $(GOPATH)/bin:$(PATH)
export AUDIT_STUB_ENABLED=1
export AUDIT_FILE=./audit.log
export AUDIT_URL=http://localhost:8000
export DEBUG_ADDR=localhost:6060

.PHONY: all tidy build run test mock migrate-up migrate-down clean

all: build

## ---------------------------
## Dependencies
## ---------------------------

tidy:
	$(GO) mod tidy

tidyvendor:
	$(GO) mod tidy
	$(GO) mod vendor

## ---------------------------
## Build & Run
## ---------------------------

build:	
	$(GO) build -o bin/$(BINARY) \
		-ldflags "-X 'main.buildVersion=$(BUILD_VER)' \
		          -X 'main.buildDate=$(BUILD_DATE)' \
		          -X 'main.buildCommit=$(BUILD_COMMIT)'" \
		$(SHORTENER_PKG)
run:
	$(GO) run $(PKG) -d $(DB_URL) --config ./config/config.json

clean:
	rm -rf bin

## ---------------------------
## Testing
## ---------------------------

test:
	$(GO) test ./... -cover -v

coverage: 
	$(GO) test ./internal/... ./pkg/... -covermode=count -coverpkg=./... -coverprofile=coverage.out
	grep -vE "(mock\.go|/mock/|/*.gen.go|/testutils/)" coverage.out > coverage.filtered.out
	mv coverage.filtered.out coverage.out
	$(GO) tool cover -func=coverage.out

lint:
	$(GO) run ./cmd/staticlint $(PKG)   

## ---------------------------
## Generate mocks and stuff
## ---------------------------

generate:
	$(GO) generate ./...

gen-reset:
	$(GO) go run ./cmd/reset/main.go $(PKG)

## ---------------------------
## Migrations
## ---------------------------

migrate-up:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) up

migrate-down:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) down

migrate-new:
	$(GOPATH)/bin/migrate -database $(DB_URL) create -ext sql -seq -digits 6 -dir $(MIGRATIONS_DIR) $(NAME)

migrate-force:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) force $(VER)