#!/usr/bin/env bash

set -euo pipefail

test_root="$(mktemp -d)"
pids=()
app_port="${TEST_APP_PORT:-4173}"
base_path="${TEST_BASE_PATH:-/}"
multi_instance="${TEST_MULTI_INSTANCE:-false}"

cleanup() {
  for pid in "${pids[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  for pid in "${pids[@]}"; do
    wait "$pid" 2>/dev/null || true
  done
  rm -rf "$test_root"
}
trap cleanup EXIT INT TERM

bun run build -- --outDir "$test_root/dashboard" --emptyOutDir

cd ../server
go build -o "$test_root/lovely-eye-server" ./cmd/server

start_instance() {
  local instance_base_path="$1"
  local instance_port="$2"
  local database_name="$3"
  env \
    SERVER_HOST=127.0.0.1 \
    SERVER_PORT="$instance_port" \
    BASE_PATH="$instance_base_path" \
    DASHBOARD_PATH="$test_root/dashboard" \
    DB_DRIVER=sqlite \
    DB_DSN="file:$test_root/$database_name.db?cache=shared&mode=rwc" \
    JWT_SECRET=playwright-test-secret-with-at-least-32-characters \
    INITIAL_ADMIN_USERNAME=e2e-admin \
    INITIAL_ADMIN_PASSWORD=e2e-password \
    ALLOW_REGISTRATION=false \
    SECURE_COOKIES=false \
    AUTH_RATE_LIMIT_ENABLED=false \
    ANALYTICS_RATE_LIMIT_ENABLED=false \
    GEOIP_DB_PATH="$test_root/$database_name.mmdb" \
    GEOIP_DOWNLOAD_URL=http://127.0.0.1:1/unavailable \
    "$test_root/lovely-eye-server" &
  pids+=("$!")
}

if [[ "$multi_instance" == "true" ]]; then
  go build -o "$test_root/browser-test-proxy" ./e2e/browserproxy
  start_instance /instance-a 4174 instance-a
  start_instance /instance-b 4175 instance-b
  TEST_APP_PORT="$app_port" "$test_root/browser-test-proxy" &
  proxy_pid=$!
  pids+=("$proxy_pid")
  wait "$proxy_pid"
else
  start_instance "$base_path" "$app_port" lovely-eye
  wait "${pids[0]}"
fi
