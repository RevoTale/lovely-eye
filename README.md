# Lovely Eye

Privacy-focused web analytics. A simple, blazingly fast, self-hosted alternative to Google Analytics, Umami and Plausible. Designed to work well with low resource systems.
Written in Go programming language. 
![Lovely Eye Logo Banner](./preview.png "Lovely Eye")
> [!WARNING]
> Work in progress

## Features

- **Privacy-first**: No cookies, daily visitor ID rotation.
- **Bot filtering**: Excludes crawlers, scrapers, and monitoring bots
- **Lightweight**: Low-RAM Docker builds, SQLite or PostgreSQL for data persitance. Lightweight tracking script
- **Real-time dashboard**: GraphQL API with modern React UI
- **Custom events**: Track clicks and user interactions via custom events

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
      # Optional: enable country stats with a MaxMind license key (auto-downloads to /data)
      - GEOIP_MAXMIND_LICENSE_KEY=your-maxmind-license-key
    volumes:
      - lovely-eye-data:/app/data
      # Optional: mount /data once for both SQLite and GeoIP
      - ./data:/data
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
      # Optional: enable country stats with a MaxMind license key (auto-downloads to /data)
      - GEOIP_MAXMIND_LICENSE_KEY=your-maxmind-license-key
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    volumes:
      # Optional: mount /data once for GeoIP files
      - ./data:/data

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

Server starts at http://localhost:8080. The first registered user becomes admin. SQLite by default.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `DB_DRIVER` | `sqlite` | `sqlite` or `postgres` |
| `DB_DSN` | `file:data/lovely_eye.db?cache=shared&mode=rwc` | Database connection string |
| `JWT_SECRET` | (random) | JWT signing key (min 32 chars, required for production) |
| `SECURE_COOKIES` | `true` | Use secure cookies (requires HTTPS). Set to `false` for local dev |
| `ALLOW_REGISTRATION` | `false` | Allow new user registration after first user |
| `GEOIP_DB_PATH` | `/data/GeoLite2-Country.mmdb` | Path to GeoLite2-Country.mmdb for country stats |
| `GEOIP_DOWNLOAD_URL` | `https://download.db-ip.com/free/dbip-country-lite-YYYY-MM.mmdb.gz` | Custom GeoIP download URL (mmdb, gz, or tar.gz). DB-IP URLs will try the current and previous 2 monthly filenames automatically. |
| `GEOIP_MAXMIND_LICENSE_KEY` | - | MaxMind license key for GeoLite2 auto-download |

Country tracking downloads the GeoLite2 database on demand when at least one site enables it. If the download fails, the dashboard will show the error in site settings.

## Authentication

JWT tokens in HttpOnly cookies with SameSite settings:

- **HttpOnly**: Prevents JavaScript access (XSS protection)
- **Secure**: HTTPS only in production
- **SameSite=Strict** (production) or **Lax** (development): Prevents CSRF

No CSRF tokens needed. See [discussion](https://www.reddit.com/r/node/comments/1im7yj0/comment/mc0ylfd/).

## How tracking work
For main mechanics of the tracking look at the [ANALYTICS.md](./ANALYTICS.md)

## How do we handle privacy
To learn more about how we handle the privacy take a look at the [PRIVACY.md](./PRIVACY.md).

### Track Custom Events

```html
<script>
  // Example: track an error event with metadata
  window.lovelyEye?.track('error', {
    message: 'Checkout failed',
    code: 'PAYMENT_DECLINED',
  });
</script>
```

Events must be allowlisted in the site settings. Unknown event names or fields are ignored.

## License

Copyright 2025 RevoTale

Licensed under [AGPL-3.0](./LICENSE).
