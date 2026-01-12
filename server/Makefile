.PHONY: build run test clean generate migrate-init migrate-up migrate-down migrate-status migrator-diff migrator-hash test-migrations test-migrations-sqlite test-migrations-postgres

build:
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate
	go build -o bin/migrator ./cmd/migrator

run:
	go run ./cmd/server

test:
	go test -v ./... --count=2

clean:
	rm -rf bin/
	rm -rf data/*

generate:
	go generate ./...

# Production migration commands (runtime)
migrate-init:
	go run ./cmd/migrate init

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-status:
	go run ./cmd/migrate status

# Development migration commands (preparation)
migrator-diff:
	go run ./cmd/migrator diff

migrator-hash:
	go run ./cmd/migrator hash

# Migration testing (runs full up/down cycle)
test-migrations-sqlite:
	@../scripts/test-migrations-sqlite.sh

test-migrations-postgres:
	@../scripts/test-migrations-postgres.sh

test-migrations: test-migrations-sqlite test-migrations-postgres
