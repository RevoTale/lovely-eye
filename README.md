# Lovely Eye

![Tracker Size Badge](./server/static/dist/tracker-size.svg "Tracker size")

Self-hosted web analytics with a Go backend and a React dashboard. Lovely Eye tracks page views and allowlisted custom events without analytics cookies or client-side identifiers. It runs on SQLite by default, supports PostgreSQL, and is designed to stay lightweight on small hosts.

## Highlights

- Cookieless analytics with an identifier computed from minimized request data and keyed with a server-side secret
- SQLite by default, PostgreSQL when needed
- Bot filtering and page-view deduplication
- Allowlisted custom events
- Optional country tracking
- Dashboard served as static assets by the Go server
- Extremely lighweight runtime

## Quick Start

The Docker Compose examples below are meant to be copied directly. They use `SECURE_COOKIES=false` so dashboard auth works on `http://localhost`. Change it to `true` when you serve Lovely Eye behind HTTPS.

### Docker Compose (SQLite)

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:latest
    ports:
      - "8080:8080"
    environment:
      - JWT_SECRET=replace-with-a-32-plus-character-secret
      - ANALYTICS_IDENTITY_SECRET=replace-with-a-second-32-plus-character-secret
      - SECURE_COOKIES=false
      # Leave both empty to allow the first registered user to become admin.
      # Set both to create the initial admin on startup.
      - INITIAL_ADMIN_USERNAME=
      - INITIAL_ADMIN_PASSWORD=
    volumes:
      - lovely-eye-data:/app/data
      - ./data:/data
    restart: unless-stopped

volumes:
  lovely-eye-data:
```

```bash
docker compose up -d
```

Open `http://localhost:8080/dashboard`.

### Docker Compose (PostgreSQL)

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:latest
    ports:
      - "8080:8080"
    environment:
      - DB_DRIVER=postgres
      - DB_DSN=postgres://lovely:lovely@lovely-eye-db:5432/lovely_eye?sslmode=disable
      - JWT_SECRET=replace-with-a-32-plus-character-secret
      - ANALYTICS_IDENTITY_SECRET=replace-with-a-second-32-plus-character-secret
      - SECURE_COOKIES=false
      - INITIAL_ADMIN_USERNAME=
      - INITIAL_ADMIN_PASSWORD=
    depends_on:
      lovely-eye-db:
        condition: service_healthy
    restart: unless-stopped

  lovely-eye-db:
    image: postgres:18.3-alpine
    environment:
      - POSTGRES_USER=lovely
      - POSTGRES_PASSWORD=lovely
      - POSTGRES_DB=lovely_eye
    volumes:
      - lovely-eye-postgres:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U lovely -d lovely_eye"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  lovely-eye-postgres:
```

```bash
docker compose up -d
```

### From Source

Requires Go 1.26+.

```bash
cd server
go run ./cmd/server
```

SQLite is the default database. If `JWT_SECRET` is unset, Lovely Eye generates one at startup and dashboard sessions do not survive restarts. If `ANALYTICS_IDENTITY_SECRET` is unset, analytics falls back to `JWT_SECRET`.

## Initial Admin And Registration

- If both `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` are set, Lovely Eye creates that initial admin on startup.
- If both initial-admin values are set and `ALLOW_REGISTRATION` is unset or empty, post-bootstrap registration defaults to disabled.
- If either initial-admin value is missing, no initial admin is created.
- If either initial-admin value is missing and `ALLOW_REGISTRATION` is unset or empty, registration defaults to enabled.
- The first registration is always available when the database has no users.
- A non-empty `ALLOW_REGISTRATION` value explicitly overrides the derived default.

## Privacy And Tracking

- Lovely Eye does not use analytics cookies or local storage by default.
- The analytics visitor identifier is computed from site ID, truncated IP prefix, browser family, and device class, and keyed with a server-side secret.
- The analytics visitor identifier is unique per site.
- The server computes hashes for `today` and `yesterday`.
- A visitor who returns at least once per UTC day keeps the same analytics client row.
- A new analytics client row is created only after the visitor skips a full UTC day between visits.
- Sessions are separate from the analytics visitor identifier and expire after 30 minutes of inactivity.
- Country tracking is optional and is not part of the analytics visitor identifier.
- The dedicated `ANALYTICS_IDENTITY_SECRET` helps reduce the impact of database-only leaks because stored analytics rows do not contain enough information to recompute the identifier on their own.

## Install The Tracker

1. Sign in to the dashboard.
2. Create a site.
3. Open the site settings.
4. Copy the generated tracking code.
5. Add it to the site you want to track.

## Common Configuration

| Variable | Default | Meaning |
| --- | --- | --- |
| `DB_DRIVER` | `sqlite` | Database driver: `sqlite` or `postgres` |
| `DB_DSN` | `file:data/lovely_eye.db?cache=shared&mode=rwc` | Database connection string |
| `JWT_SECRET` | generated at startup if empty | Dashboard auth secret. Set it explicitly in production. |
| `ANALYTICS_IDENTITY_SECRET` | falls back to `JWT_SECRET` | Optional dedicated secret for the analytics visitor identifier |
| `SECURE_COOKIES` | `true` | Enables secure dashboard auth cookies |
| `ALLOW_REGISTRATION` | `auto` | Empty means derived from the initial-admin envs |
| `INITIAL_ADMIN_USERNAME` | empty | Initial admin username. Requires `INITIAL_ADMIN_PASSWORD`. |
| `INITIAL_ADMIN_PASSWORD` | empty | Initial admin password. Requires `INITIAL_ADMIN_USERNAME`. |
| `GEOIP_MAXMIND_LICENSE_KEY` | empty | Optional MaxMind license key for country tracking |

## Custom Events

```html
<script>
  window.lovelyEye?.track("checkout_failed", {
    code: "PAYMENT_DECLINED",
    step: "confirm",
  });
</script>
```

Custom events are recorded only when the event name and fields are allowlisted in site settings. Otherwise, they are discarded silently.

## Documentation

- [ANALYTICS.md](./ANALYTICS.md) - tracking mechanics
- [PRIVACY.md](./PRIVACY.md) - privacy handling
- [dashboard/README.md](./dashboard/README.md) - dashboard development
- [server/CONTRIBUTING.md](./server/CONTRIBUTING.md) - server development notes

## Advanced Docker Compose Example

This example includes all server environment variables. Start with the quick-start examples unless you need to tune these values explicitly.

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:latest
    ports:
      - "8080:8080"
    environment:
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
      - BASE_PATH=/
      - DASHBOARD_PATH=dashboard
      - DB_DRIVER=postgres
      - DB_DSN=postgres://lovely:lovely@lovely-eye-db:5432/lovely_eye?sslmode=disable
      - DB_MAX_CONNS=10
      - DB_MIN_CONNS=1
      - DB_CONNECT_TIMEOUT=7s
      - JWT_SECRET=replace-with-a-32-plus-character-secret
      - JWT_ACCESS_EXPIRY_MINUTES=15
      - JWT_REFRESH_DAYS=7
      - SECURE_COOKIES=true
      - COOKIE_DOMAIN=
      # Leave empty for the derived default:
      # false when both INITIAL_ADMIN_* values are set, true otherwise.
      - ALLOW_REGISTRATION=
      # Set both or leave both empty.
      - INITIAL_ADMIN_USERNAME=
      - INITIAL_ADMIN_PASSWORD=
      - ANALYTICS_IDENTITY_SECRET=replace-with-a-second-32-plus-character-secret
      - GEOIP_DB_PATH=/data/GeoLite2-Country.mmdb
      - GEOIP_DOWNLOAD_URL=https://download.db-ip.com/free/dbip-country-lite.mmdb.gz
      - GEOIP_MAXMIND_LICENSE_KEY=
      - LOG_LEVEL=warn
    depends_on:
      lovely-eye-db:
        condition: service_healthy
    volumes:
      - lovely-eye-data:/data
    restart: unless-stopped

  lovely-eye-db:
    image: postgres:18.3-alpine
    environment:
      - POSTGRES_USER=lovely
      - POSTGRES_PASSWORD=lovely
      - POSTGRES_DB=lovely_eye
    volumes:
      - lovely-eye-postgres:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U lovely -d lovely_eye"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  lovely-eye-data:
  lovely-eye-postgres:
```
## Banner
![Lovely Eye Logo Banner](./preview.png "Lovely Eye")

## License

Licensed under [AGPL-3.0-or-later](./LICENSE). See [COPYRIGHT](./COPYRIGHT).
