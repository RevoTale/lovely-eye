#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "Testing SQLite migrations..."
docker compose -f docker/docker-compose.migrations-test.yml --profile sqlite run --rm --build test-migrations-sqlite
EXIT_CODE=$?

echo "Cleaning up..."
docker compose -f docker/docker-compose.migrations-test.yml --profile sqlite down -v

exit $EXIT_CODE
