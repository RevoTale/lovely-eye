.PHONY: build run test clean generate migrate-init migrate-up migrate-down migrate-status migrator-diff migrator-hash test-migrations test-migrations-sqlite test-migrations-postgres

build:
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate

run:
	go run ./cmd/server

test:
	go test ./... --count=2

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

# Development migration commands (preparation using Atlas CLI)
migrator-diff:
	@read -p "Migration name: " name; \
	echo "\nGenerating SQLite migration..."; \
	atlas migrate diff $$name --env sqlite; \
	echo "\nGenerating PostgreSQL migration..."; \
	atlas migrate diff $$name --env postgres; \
	echo "\nMigrations generated successfully!"

migrator-hash:
	@echo "Hashing SQLite migrations..."; \
	atlas migrate hash --env sqlite; \
	echo "Hashing PostgreSQL migrations..."; \
	atlas migrate hash --env postgres; \
	echo "\nHash files updated successfully!"

# Migration testing (runs full up/down cycle)
test-migrations-sqlite:
	@../scripts/test-migrations-sqlite.sh

test-migrations-postgres:
	@../scripts/test-migrations-postgres.sh

test-migrations: test-migrations-sqlite test-migrations-postgres
