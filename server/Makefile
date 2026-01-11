.PHONY: build run test clean generate migrate migrate-down migrate-status migrate-init atlas-diff atlas-hash

build:
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate

run:
	go run ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -f lovely_eye.db

generate:
	go generate ./...

# Migration commands
migrate-init:
	go run ./cmd/migrate db init

migrate:
	go run ./cmd/migrate db migrate

migrate-down:
	go run ./cmd/migrate db rollback

migrate-status:
	go run ./cmd/migrate db status

# Atlas - auto-generate migrations for both databases
atlas-diff:
	@read -p "Migration name: " name; \
	atlas migrate diff $$name --env sqlite && \
	atlas migrate diff $$name --env postgres --dev-url "sqlite://file?mode=memory"

atlas-hash:
	@atlas migrate hash --env sqlite
	@atlas migrate hash --env postgres
