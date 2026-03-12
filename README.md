# Lovely Eye

![Tracker Size Badge](./server/static/dist/tracker-size.svg "Tracker size")

Self-hosted web analytics with a Go backend and React dashboard. Built for low-resource hosts, with Go's small memory footprint and single-binary deployment keeping it lightweight. Supports SQLite or PostgreSQL.


## Features

- **Privacy-first**: no analytics cookies, keyed server-side visitor identity with UTC-day-skipped rotation
- **Bot filtering**: excludes crawlers, scrapers, monitoring bots.
- **Lightweight**: runtime consumes around ~15MB of RAM on AMD processor.
- **SQLite and PostgreSQL** supported.
- **Dashboard**: GraphQL API with React UI deployed as a static assets.
- **Custom events**: allowlisted event names and fields. Prevents spamming your database with unnecessary data.

## Quick Start

### Docker (SQLite)

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:latest
    ports:
      - "8080:8080"
    environment:
      # Optional for local development. Set a fixed value in production to
      # preserve dashboard sessions across restarts.
      - JWT_SECRET=${JWT_SECRET:-}
      # Optional: dedicated secret for analytics visitor identity. Falls back to
      # JWT_SECRET when unset. Set a fixed value in production to keep visitor
      # counting stable across restarts.
      - ANALYTICS_IDENTITY_SECRET=${ANALYTICS_IDENTITY_SECRET:-}
      # Dashboard auth uses cookies and requires HTTPS (serve behind a reverse proxy)
      # Default is true; you can disable it to serve over HTTP.
      # Tracking is cookieless; this setting only affects dashboard auth.
      - SECURE_COOKIES=true
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
      - "${PORT:-8080}:8080"
    environment:
      - DB_DRIVER=postgres
      - DB_DSN=postgres://${POSTGRES_USER:-lovely}:${POSTGRES_PASSWORD:-lovely}@lovely-eye-db:5432/${POSTGRES_DB:-lovely_eye}?sslmode=disable
      # Optional for local development. Set a fixed value in production to
      # preserve dashboard sessions across restarts.
      - JWT_SECRET=${JWT_SECRET:-}
      # Optional: dedicated secret for analytics visitor identity. Falls back to
      # JWT_SECRET when unset. Set a fixed value in production to keep visitor
      # counting stable across restarts.
      - ANALYTICS_IDENTITY_SECRET=${ANALYTICS_IDENTITY_SECRET:-}
      # Dashboard auth uses cookies and requires HTTPS (serve behind a reverse proxy)
      # Default is true; you can disable it to serve over HTTP.
      # Tracking is cookieless; this setting only affects dashboard auth.
      - SECURE_COOKIES=true
      - INITIAL_ADMIN_PASSWORD=${INITIAL_ADMIN_PASSWORD}
      - INITIAL_ADMIN_USERNAME=${INITIAL_ADMIN_USERNAME}
    depends_on:
      lovely-eye-db:
        condition: service_healthy
    networks:
      - lovely-eye-net
    restart: unless-stopped

  lovely-eye-db:
    image: postgres:18.1-alpine
    environment:
      - POSTGRES_USER=${POSTGRES_USER:-lovely}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-lovely}
      - POSTGRES_DB=${POSTGRES_DB:-lovely_eye}
    volumes:
      - lovely-eye-data:/var/lib/postgresql
    networks:
      - lovely-eye-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-lovely} -d ${POSTGRES_DB:-lovely_eye}"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  lovely-eye-data:

networks:
  lovely-eye-net:

```


### From Source

Requires Go 1.26+.

```bash
cd server
go run ./cmd/server
```
Server starts at http://localhost:8080. SQLite by default.
If `JWT_SECRET` is unset, the server generates one at startup and dashboard sessions do not survive restarts.
If `ANALYTICS_IDENTITY_SECRET` is unset, analytics identity falls back to `JWT_SECRET`.

Registration defaults depend on initial-admin config. If both `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` are set, Lovely Eye creates that admin on startup and keeps registration disabled unless `ALLOW_REGISTRATION=true` is explicitly set. If either initial-admin value is missing, the first self-registered user becomes admin and registration stays enabled by default.

### Next Step. Install the tracking script.

After you started your containers:
- login into dashboard.
- Create "Site" and enter domain you want to track.
- Open your "Site". You will find the "Settings" button in right top corner. Click it and scroll to the "Tracking Code" section.
- Copy the tracking code and install into your website via `javascript` or `html`.
<img width="1060" height="586" alt="Screenshot 2026-01-26 at 20 58 55" src="https://github.com/user-attachments/assets/e4f878dd-fca2-4678-bf49-3dddf04920a2" />



## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `DB_DRIVER` | `sqlite` | `sqlite` or `postgres` |
| `DB_DSN` | `file:data/lovely_eye.db?cache=shared&mode=rwc` | Database connection string |
| `JWT_SECRET` | generated at startup if empty | JWT signing key. Must be at least 32 chars when set. Set it explicitly in production or multi-instance deployments. |
| `ANALYTICS_IDENTITY_SECRET` | falls back to `JWT_SECRET` | Optional dedicated secret for analytics visitor identity. Must be at least 32 chars when set. Helps reduce the impact of database-only leaks by making visitor IDs harder to recompute. |
| `SECURE_COOKIES` | `true` | Use secure cookies (requires HTTPS). Set to `false` for local dev |
| `ALLOW_REGISTRATION` | `auto` | Post-bootstrap registration policy. Defaults to `false` when both `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` are set, otherwise defaults to `true`. The first registration is still available whenever no users exist. |
| `INITIAL_ADMIN_USERNAME` | (empty) | Optional initial admin username. Takes effect only when `INITIAL_ADMIN_PASSWORD` is also set. |
| `INITIAL_ADMIN_PASSWORD` | (empty) | Optional initial admin password. Takes effect only when `INITIAL_ADMIN_USERNAME` is also set. |
| `GEOIP_DB_PATH` | `/data/GeoLite2-Country.mmdb` | Path to GeoLite2-Country.mmdb for country stats |
| `GEOIP_DOWNLOAD_URL` | `https://download.db-ip.com/free/dbip-country-lite.mmdb.gz` | Default GeoIP download URL (mmdb, gz, or tar.gz) |
| `GEOIP_MAXMIND_LICENSE_KEY` | - | MaxMind license key for GeoLite2 auto-download |

Country tracking downloads the GeoIP database on demand when at least one site enables it. If the download fails, the dashboard will show the error in site settings.

Analytics visitor identity is server-generated and derived from a keyed hash of site ID, truncated IP prefix, browser family, and device class. The hash is computed per UTC day, but the server reuses the same client across `today` and `yesterday`; if only yesterday matches, that row is rewritten to today's hash. Country tracking stays separate from visitor identity, sessions still expire after 30 minutes of inactivity, and the dedicated analytics identity secret helps reduce the impact of database-only leaks because visitor IDs cannot be recomputed from stored analytics data alone.

## Custom Events

```html
<script>
  window.lovelyEye?.track('error', {
    message: 'Checkout failed',
    code: 'PAYMENT_DECLINED',
  });
</script>
```

Events must be allowlisted in site settings. Unknown event names or fields are ignored.

## Showcase screenshots
<img width="1512" height="982" alt="Screenshot 2026-01-26 at 15 21 10" src="https://github.com/user-attachments/assets/a231dad7-02dc-442d-8d7c-2d7dd459c05d" />
<img width="1512" height="885" alt="Screenshot 2026-01-26 at 15 23 44" src="https://github.com/user-attachments/assets/f7916173-9a92-4502-b055-aca27089205d" />



## Documentation

- [ANALYTICS.md](./ANALYTICS.md) - tracking mechanics
- [PRIVACY.md](./PRIVACY.md) - privacy handling

## Authentication

JWT tokens in HttpOnly cookies with SameSite settings:

- **HttpOnly**: prevents JavaScript access (XSS protection)
- **Secure**: HTTPS only in production
- **SameSite=Strict** (production) or **Lax** (development): prevents CSRF


## License

Copyright 2025 RevoTale

Licensed under [AGPL-3.0](./LICENSE).

![Lovely Eye Logo Banner](./preview.png "Lovely Eye")
