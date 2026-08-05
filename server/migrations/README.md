# Database Migrations

Atlas auto-generates SQL migrations from Bun models, Bun applies them at runtime.

## Prerequisites

The devcontainer automatically provides:
- Atlas CLI (installed in Dockerfile)
- PostgreSQL 18 (via docker-compose service)

## Development Workflow

### Create New Migrations
1. Edit models in `internal/models/models.go`
2. Run `task migrator-diff` (prompts for migration name)
3. Atlas CLI generates `.up.sql` and `.down.sql` for both SQLite and PostgreSQL

### Test Before Committing
```bash
task test-migrations  # Tests full up/down cycle on both databases
```

### Apply Locally
```bash
task migrate-up
```

## Production Deployment

```bash
task migrate-init    # First time only - creates migration tracking tables
task migrate-up      # Applies all pending migrations
task migrate-status  # Shows what's applied
task migrate-down    # Rollback if needed
```

## CI/CD Integration

Two parallel test jobs run on every release:
- `test-migrations-sqlite` - Tests on SQLite
- `test-migrations-postgres` - Tests on real PostgreSQL

Both must pass before Docker images are published. Each test:
1. Builds with production Dockerfile
2. Runs all migrations up
3. Rolls all migrations down
4. Applies them again (idempotency check)

## Structure

```
migrations/
├── sqlite/      # SQLite-specific migrations
├── postgres/    # PostgreSQL-specific migrations
└── atlas-schema.go  # Schema definition for Atlas
```

Separate directories needed because SQLite and PostgreSQL use different syntax (e.g., `INTEGER AUTOINCREMENT` vs `BIGSERIAL`).

## Environment Variables

- `DB_DRIVER` - `sqlite` (default) or `postgres`
- `DB_DSN` - Connection string
- `JWT_SECRET` - Optional. If unset, the app generates one at startup. Set it explicitly in production because dashboard sessions will not survive restarts.
- `AUTH_RATE_LIMIT_ENABLED` - Optional. Defaults to `true` for `login` and `register` GraphQL mutations.
- `AUTH_RATE_LIMIT_ATTEMPTS` - Optional. Defaults to `10` per trusted client IP during the window.
- `AUTH_RATE_LIMIT_WINDOW` - Optional. Defaults to `15m`.
- `ANALYTICS_IDENTITY_SECRET` - Optional. Falls back to `JWT_SECRET`. Set it explicitly in production if visitor identity should remain stable across restarts without sharing the auth secret. Analytics uses it for the daily UTC hashes behind UTC-day-skipped rotation, and it also reduces the impact of database-only leaks by making visitor IDs harder to recompute.
- `ANALYTICS_MAX_BODY_BYTES` - Optional. Defaults to `16384` for small collect payloads.
- `ANALYTICS_MAX_PROPERTIES_BYTES` - Optional. Defaults to `8192` for custom event properties.
- `ANALYTICS_RATE_LIMIT_ENABLED` - Optional. Defaults to `true` for the public collect endpoint.
- `ANALYTICS_RATE_LIMIT_PER_MINUTE` - Optional. Defaults to `120` per site key and client IP.
- `ANALYTICS_RATE_LIMIT_BURST` - Optional. Defaults to `240` per site key and client IP.
- `TRUSTED_PROXY_CIDRS` - Optional. Defaults to loopback, private IPv4 ranges, and IPv6 unique-local addresses. Public CDN ranges must be configured explicitly.
- `GRAPHQL_MAX_BODY_BYTES` - Optional. Defaults to `1048576`.
- `DASHBOARD_MAX_DAILY_RANGE_DAYS` - Optional. Defaults to `730`.
- `DASHBOARD_MAX_HOURLY_RANGE_DAYS` - Optional. Defaults to `31`.
- `DASHBOARD_MAX_FILTER_VALUES` - Optional. Defaults to `100`.
- `DASHBOARD_MAX_FILTER_STRING_LENGTH` - Optional. Defaults to `2048`.
