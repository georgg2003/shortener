install:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go get -u github.com/golang-migrate/migrate/v4

migrate:
	migrate -database "postgres://shortener:password@localhost:5432/shortener?sslmode=disable" -path ./migrations up

run:
	go run ./...