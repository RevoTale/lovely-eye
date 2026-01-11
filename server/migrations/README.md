# Migrations

Atlas generates SQL from Bun models → Bun applies at runtime.

## Workflow

1. Edit `internal/models/models.go`
2. `make atlas-diff` - generates SQLite + PostgreSQL migrations
3. `make migrate` - applies migrations

## Runtime Database

```bash
make migrate              # SQLite
DB_DRIVER=postgres make migrate  # PostgreSQL
```

Correct migrations load automatically based on `DB_DRIVER`.

## Why Separate Files?

SQLite uses `INTEGER AUTOINCREMENT`, PostgreSQL uses `BIGSERIAL`.  
Atlas generates both from the same Bun models.
