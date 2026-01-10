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

### Docker (SQLite)

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:latest
    ports:
      - "8080:8080"
    environment:
      - JWT_SECRET=your-secret-key-min-32-chars
      - SECURE_COOKIES=true
    volumes:
      - lovely-eye-data:/app/data
    restart: unless-stopped

volumes:
  lovely-eye-data:
```

### Docker (PostgreSQL)

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:latest
    ports:
      - "8080:8080"
    environment:
      - DB_DRIVER=postgres
      - DB_DSN=postgres://lovely:lovely@postgres:5432/lovely_eye?sslmode=disable
      - JWT_SECRET=your-secret-key-min-32-chars
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  postgres:
    image: postgres:18.1-alpine
    environment:
      - POSTGRES_USER=lovely
      - POSTGRES_PASSWORD=lovely
      - POSTGRES_DB=lovely_eye
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U lovely -d lovely_eye"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  postgres-data:
```

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
| `SECURE_COOKIES` | `true` | Use secure cookies (requires HTTPS). Set to `false` for local dev |
| `ALLOW_REGISTRATION` | `false` | Allow new user registration after first user |

## Authentication

JWT tokens in HttpOnly cookies with SameSite settings:

- **HttpOnly**: Prevents JavaScript access (XSS protection)
- **Secure**: HTTPS only in production
- **SameSite=Strict** (production) or **Lax** (development): Prevents CSRF

No CSRF tokens needed. See [discussion](https://www.reddit.com/r/node/comments/1im7yj0/comment/mc0ylfd/).

## License

Copyright 2025 RevoTale

Licensed under [AGPL-3.0](./LICENSE).
