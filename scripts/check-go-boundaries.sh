#!/usr/bin/env bash
set -euo pipefail

repository_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

for legacy_root in \
  "$repository_root/server/internal/config" \
  "$repository_root/server/internal/clientip" \
  "$repository_root/server/internal/database" \
  "$repository_root/server/internal/handlers" \
  "$repository_root/server/internal/middleware" \
  "$repository_root/server/internal/models" \
  "$repository_root/server/internal/repository" \
  "$repository_root/server/internal/server" \
  "$repository_root/server/internal/services" \
  "$repository_root/server/pkg"; do
  if [[ -d "$legacy_root" ]] && find "$legacy_root" -type f -print -quit | grep -q .; then
    echo "legacy backend root still contains files: ${legacy_root#"$repository_root/"}" >&2
    exit 1
  fi
done

if grep -RInE \
  'github.com/lovely-eye/server/(internal/(clientip|config|database|handlers|middleware|models|repository|server|services)|pkg/)' \
  "$repository_root/server" \
  --include='*.go'; then
  echo "backend code imports a legacy package path" >&2
  exit 1
fi

if grep -RInE \
  'github.com/lovely-eye/server/internal/.+/persistence' \
  "$repository_root/server/internal/graph" \
  "$repository_root/server/internal/transport" \
  --include='*.go' \
  --exclude='*_test.go'; then
  echo "transport code must consume feature contracts, not persistence adapters" >&2
  exit 1
fi

for persistence_root in "$repository_root"/server/internal/*/persistence; do
  if grep -RInE \
    'github.com/lovely-eye/server/internal/(app|graph|seed|transport)' \
    "$persistence_root" \
    --include='*.go' \
    --exclude='*_test.go'; then
    echo "persistence adapters may not depend on application or transport packages" >&2
    exit 1
  fi
done

if grep -RInE \
  'github.com/lovely-eye/server/internal/(analytics|auth|country|event|geoip|graph|handlers|middleware|models|repository|site|transport)' \
  "$repository_root/server/cmd" \
  --include='*.go'; then
  echo "cmd packages must delegate to app, seed, migration, or platform entry points" >&2
  exit 1
fi

echo "Go architecture boundaries are valid."
