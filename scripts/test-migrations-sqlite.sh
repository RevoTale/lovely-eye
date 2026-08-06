#!/bin/bash
set -euo pipefail

repository_root="$(cd "$(dirname "$0")/.." && pwd)"
compose_file="$repository_root/docker/docker-compose.migrations-test.yml"

cleanup() {
  docker compose -f "$compose_file" --profile sqlite down -v
}
trap cleanup EXIT

echo "Testing SQLite migrations..."
docker compose -f "$compose_file" --profile sqlite run --rm --build test-migrations-sqlite
