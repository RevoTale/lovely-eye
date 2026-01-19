# Database Migrations

Atlas auto-generates SQL migrations from Bun models, Bun applies them at runtime.

## Prerequisites

The devcontainer automatically provides:
- Atlas CLI (installed in Dockerfile)
- PostgreSQL 18 (via docker-compose service)

## Development Workflow

### Create New Migrations
1. Edit models in `internal/models/models.go`
2. Run `make migrator-diff` (prompts for migration name)
3. Atlas CLI generates `.up.sql` and `.down.sql` for both SQLite and PostgreSQL

### Test Before Committing
```bash
make test-migrations  # Tests full up/down cycle on both databases
```

### Apply Locally
```bash
make migrate-up
```

## Production Deployment

```bash
make migrate-init    # First time only - creates migration tracking tables
make migrate-up      # Applies all pending migrations
make migrate-status  # Shows what's applied
make migrate-down    # Rollback if needed
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
- `JWT_SECRET` - Required for app initialization
