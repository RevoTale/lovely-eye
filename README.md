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
| `SECURE_COOKIES` | `true` | Use secure cookies (requires HTTPS). Set to `false` for local dev |
| `ALLOW_REGISTRATION` | `false` | Allow new user registration after first user |

## Authentication

Lovely Eye uses modern, secure authentication with JWT tokens stored in HttpOnly cookies with proper SameSite settings:

- **HttpOnly**: Prevents JavaScript access to tokens (XSS protection)
- **Secure**: Cookies only sent over HTTPS in production
- **SameSite=Strict** (production) or **SameSite=Lax** (development): Prevents CSRF attacks

**No CSRF tokens needed!** Modern browsers with proper cookie settings eliminate the need for CSRF tokens. As explained in [this Reddit discussion](https://www.reddit.com/r/node/comments/1im7yj0/comment/mc0ylfd/):

> "A JWT in a HTTP-Only Secure cookie + SameSite=Strict (or Lax) is basically what you need."

This approach is simpler and more secure.

## License

Copyright 2025 RevoTale

Licensed under [AGPL-3.0](./LICENSE).
