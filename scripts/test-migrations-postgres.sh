#!/bin/bash
set -euo pipefail

repository_root="$(cd "$(dirname "$0")/.." && pwd)"
compose_file="$repository_root/docker/docker-compose.migrations-test.yml"

cleanup() {
  docker compose -f "$compose_file" --profile postgres down -v
}
trap cleanup EXIT

echo "Testing PostgreSQL migrations..."
docker compose -f "$compose_file" up -d postgres --wait
docker compose -f "$compose_file" --profile postgres run --rm --build test-migrations-postgres
