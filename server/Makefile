.PHONY: build run test clean generate migrate migrate-down migrate-status

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

migrate:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-status:
	go run ./cmd/migrate status
