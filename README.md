# Lovely Eye

Privacy-focused web analytics. A simple self-hosted alternative to Google Analytics, Umami and Plausible.

> [!WARNING]
> Work in progress

## Features

- Privacy-first: No cookies, no personal data collection
- Lightweight tracking script
- Real-time dashboard via GraphQL
- SQLite (default) or PostgreSQL
- Docker support

## Quick Start

### Docker

```bash
cd docker
cp .env.example .env
# Edit .env and set JWT_SECRET
docker compose up -d
```

For PostgreSQL, use `docker compose -f docker-compose.postgres.yml up -d`.

### From Source

Requires Go 1.25+.

```bash
cd server
go run ./cmd/server
```

Server starts at http://localhost:8080. The first registered user becomes admin.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `DB_DRIVER` | `sqlite` | `sqlite` or `postgres` |
| `DB_DSN` | `file:lovely_eye.db?cache=shared&mode=rwc` | Database connection string |
| `JWT_SECRET` | (random) | JWT signing key (min 32 chars, required for production) |
| `ALLOW_REGISTRATION` | `false` | Allow new user registration after first user |

## License

AGPL-3.0
