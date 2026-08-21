#!/usr/bin/env bash

set -euo pipefail

repository_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
measurement_root="$(mktemp -d)"
cleanup() {
  rm -rf "$measurement_root"
}
trap cleanup EXIT INT TERM

echo "ENVIRONMENT"
cd "$repository_root/server"
go version
go env GOOS GOARCH
printf 'node %s\n' "$(node --version)"
printf 'pnpm %s\n' "$(pnpm --version)"
uname -m

echo "BUNDLE_SIZES"
cd "$repository_root/dashboard"
pnpm run build --outDir "$measurement_root/dashboard" --emptyOutDir --manifest >/dev/null
pnpm exec tsx scripts/report-bundle-size.ts "$measurement_root/dashboard" --check

echo "GO_BINARY_SIZES"
cd "$repository_root/server"
go build -trimpath -o "$measurement_root/server" ./cmd/server
go build -trimpath -o "$measurement_root/migrate" ./cmd/migrate
stat --printf='server_bytes=%s\n' "$measurement_root/server"
stat --printf='migrate_bytes=%s\n' "$measurement_root/migrate"
server_bytes="$(stat --printf='%s' "$measurement_root/server")"
migrate_bytes="$(stat --printf='%s' "$measurement_root/migrate")"
if ((server_bytes > 29100000)); then
  printf 'server binary exceeds 29100000 bytes: %s\n' "$server_bytes" >&2
  exit 1
fi
if ((migrate_bytes > 21400000)); then
  printf 'migration binary exceeds 21400000 bytes: %s\n' "$migrate_bytes" >&2
  exit 1
fi

echo "SQLITE_BENCHMARKS"
go test -run '^$' -bench '^BenchmarkAnalyticsDashboardReads$' -benchmem -count=10 ./internal/analytics/persistence

echo "COLLECT_BENCHMARKS"
go test -run '^$' -bench '^BenchmarkAnalyticsHandlerCollectPageView$' -benchmem -count=10 ./internal/transport/http/collect

echo "BROWSER_BASELINE"
cd "$repository_root/dashboard"
TEST_PERFORMANCE=true pnpm exec playwright test tests/e2e/performance-baseline.spec.ts
