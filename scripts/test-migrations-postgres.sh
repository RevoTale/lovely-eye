#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "Testing PostgreSQL migrations..."
docker compose -f docker/docker-compose.migrations-test.yml up -d postgres --wait
docker compose -f docker/docker-compose.migrations-test.yml --profile postgres run --rm test-migrations-postgres
EXIT_CODE=$?

echo "Cleaning up..."
docker compose -f docker/docker-compose.migrations-test.yml --profile postgres down -v

exit $EXIT_CODE
