# Lovely Eye

![Tracker Size Badge](./server/static/dist/tracker-size.svg "Tracker size")

Self-hosted web analytics with a Go backend and a React dashboard. Lovely Eye tracks page views and allowlisted custom events without analytics cookies or client-side identifiers. It runs on SQLite by default, supports PostgreSQL, and is designed to stay lightweight on small hosts.

## Highlights

- Cookieless analytics with an identifier computed from minimized request data and keyed with a server-side secret
- SQLite by default, PostgreSQL when needed
- Bot filtering and page-view deduplication
- Server-side session timing; client timing fields are ignored
- Allowlisted custom events
- Optional country tracking
- Dashboard served as static assets by the Go server
- Lightweight runtime

## Architecture And Stack

- Dashboard: React + strict TypeScript, Vite/SWC, Tailwind CSS v4, shadcn/ui with Base UI,
  TanStack Router, Apollo GraphQL, Zod, and Recharts.
- Server: Go `net/http`, gqlgen, Bun ORM, SQLite or PostgreSQL, JWT/bcrypt authentication, and
  country-only MaxMind-compatible GeoIP.
- Tooling: Node.js, pnpm, Biome, Playwright, deterministic GraphQL code generation, Atlas
  migrations, Go tool directives, and `govulncheck`.

Frontend product ownership is `app` → `features` → `shared`; backend ownership is feature packages
with private persistence adapters plus edge-owned `graph`, `transport`, and `platform` packages.
See the [dashboard structure](dashboard/README.md), [server module map](server/internal/README.md),
and [architecture decisions](docs/decisions/0003-feature-oriented-go-modular-monolith.md).

## Quick Start

The Docker Compose examples below are meant to be copied directly. They follow the moving `v2`
major channel, which receives backward-compatible v2 releases but never a future v3 release. Pin an
exact version such as `v2.0.0`, or the digest published in the GitHub Release, for controlled
production upgrades. Do not automate production deployments from `latest`; it moves across major
versions. Read [UPGRADING.md](./UPGRADING.md) before replacing an existing container.

The examples use `SECURE_COOKIES=false` so dashboard auth works on `http://localhost`. Change it to
`true` when you serve Lovely Eye behind HTTPS.

### Docker Compose (SQLite)

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:v2
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

Open `http://localhost:8080/`.

### Docker Compose (PostgreSQL)

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:v2
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
    image: postgres:18.6-alpine
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

Requires Go 1.27+.

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
- The tracker sends a minimal payload. A page view sends `path`; an exit ping sends `path` plus `exit: true`.
- Timing is computed from server receive time. The client does not send `duration`, `screen_width`, or session state.
- Single-page exit duration is bounded by `ANALYTICS_MAX_SINGLE_PAGE_DURATION`, which defaults to `4h`; repeated exit pings cannot extend it past that cap.
- The analytics visitor identifier is computed from site ID, truncated IP prefix, browser family, and device class, and keyed with a server-side secret.
- The analytics visitor identifier is unique per site.
- The server computes hashes for `today` and `yesterday`.
- A visitor who returns at least once per UTC day keeps the same analytics client row.
- A new analytics client row is created only after the visitor skips a full UTC day between visits.
- Sessions are separate from the analytics visitor identifier and expire after 30 minutes of inactivity.
- Exit pings update the active session when they match the current path. If an exit ping names a different path while the session is still inside the 30-minute active window, the server counts that path as a page view before closing the session path; stale different-path exits are ignored.
- Country tracking is optional and is not part of the analytics visitor identifier.
- The dedicated `ANALYTICS_IDENTITY_SECRET` helps reduce the impact of database-only leaks because stored analytics rows do not contain enough information to recompute the identifier on their own.

## Install The Tracker

1. Sign in to the dashboard.
2. Create a site.
3. Open the site settings.
4. Copy the generated tracking code.
5. Add it to the site you want to track.

## Tracker API

The collect endpoint requires the public key in the query string:

```http
POST /api/collect?site_key=<public_key>
Content-Type: text/plain;charset=UTF-8
```

Page view:

```json
{ "path": "/pricing" }
```

Exit ping:

```json
{ "path": "/pricing", "exit": true }
```

Initial attribution, when present, is sent only on the first page view:

```json
{ "path": "/pricing", "referrer": "https://google.com", "utm_source": "google" }
```

The tracker uses `visibilitychange` with `sendBeacon`, with `pagehide` as a fallback. This follows the current MDN and W3C Beacon guidance for small analytics payloads that should not block navigation: [MDN sendBeacon](https://developer.mozilla.org/en-US/docs/Web/API/Navigator/sendBeacon), [MDN visibilitychange](https://developer.mozilla.org/en-US/docs/Web/API/Document/visibilitychange_event), and [W3C Beacon](https://www.w3.org/TR/beacon/).

## Common Configuration

Production and release constraints are maintained in [Security and privacy model](docs/security.md).

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
| `AUTH_RATE_LIMIT_ENABLED` | `true` | Enables per-process rate limiting for `login` and `register` GraphQL mutations. |
| `AUTH_RATE_LIMIT_ATTEMPTS` | `10` | Maximum auth mutation attempts per trusted client IP during the window. |
| `AUTH_RATE_LIMIT_WINDOW` | `15m` | Fixed auth rate-limit window. |
| `GEOIP_MAXMIND_LICENSE_KEY` | empty | Optional MaxMind license key for country tracking |
| `ANALYTICS_MAX_BODY_BYTES` | `16384` | Maximum collect request body size. Small because tracker payloads are tiny. |
| `ANALYTICS_MAX_PROPERTIES_BYTES` | `8192` | Maximum custom-event `properties` JSON string size. |
| `ANALYTICS_MAX_SINGLE_PAGE_DURATION` | `4h` | Maximum same-path single-page duration accepted from an exit ping. |
| `ANALYTICS_RATE_LIMIT_ENABLED` | `true` | Enables per-process collect rate limiting. |
| `ANALYTICS_RATE_LIMIT_PER_MINUTE` | `120` | Refill rate for client IP admission and validated site key plus client IP admission. |
| `ANALYTICS_RATE_LIMIT_BURST` | `240` | Short burst allowance for the same collect admission keys. |
| `TRUSTED_PROXY_CIDRS` | private, loopback, and unique-local ranges | CIDRs allowed to supply `X-Forwarded-For` / `X-Real-IP`. Public CDN ranges must be configured explicitly. |
| `GRAPHQL_MAX_BODY_BYTES` | `1048576` | Maximum GraphQL request body size. |
| `GRAPHQL_MAX_COMPLEXITY` | `300` | Maximum calculated GraphQL operation complexity. |
| `DASHBOARD_MAX_DAILY_RANGE_DAYS` | `730` | Maximum daily dashboard date range. |
| `DASHBOARD_MAX_HOURLY_RANGE_DAYS` | `31` | Maximum hourly dashboard date range. |
| `DASHBOARD_MAX_FILTER_VALUES` | `100` | Maximum values per dashboard filter list. |
| `DASHBOARD_MAX_FILTER_STRING_LENGTH` | `2048` | Maximum byte length for each dashboard filter string. |

## Custom Events

```html
<script>
  window.lovelyEye?.track({
    name: "checkout_failed",
    properties: {
      code: "PAYMENT_DECLINED",
      step: "confirm",
    },
  });
</script>
```

Custom events are recorded only when the event name and fields are allowlisted in site settings. Otherwise, they are discarded silently.

## Development Quality Gates

Run `task check` before handing off changes. It runs frontend boundaries, full Biome, strict
TypeScript, unit tests, Go tests and lint, generated-artifact freshness, and clean production builds.
Run `task test` when the complete browser matrix is also required. Use `task lint` for lint-only work.

The dashboard uses Biome's recommended rules plus a local source-size check. The Go backend uses golangci-lint's standard v2 set plus explicit security, correctness, function-length, complexity, and file-length checks.

## Documentation

- [ANALYTICS.md](./ANALYTICS.md) - tracking mechanics
- [PRIVACY.md](./PRIVACY.md) - privacy handling
- [UPGRADING.md](./UPGRADING.md) - container backup, upgrade, verification, and rollback
- [dashboard/README.md](./dashboard/README.md) - dashboard development
- [Performance runbook](./docs/performance.md) - server RAM, allocation, CPU, and profiling baselines
- [Release runbook](./docs/releasing.md) - maintainer release contract and checklist
- [Architecture normalization plan](./docs/plans/architecture-normalization.md) - living intent, decisions, clarifications, and phased roadmap
- [ADR-0001](./docs/decisions/0001-runtime-base-path.md) - arbitrary runtime base-path support
- [server/CONTRIBUTING.md](./server/CONTRIBUTING.md) - server development notes

## Advanced Docker Compose Example

This example includes all server environment variables. Start with the quick-start examples unless you need to tune these values explicitly.

```yaml
services:
  lovely-eye:
    image: ghcr.io/revotale/lovely-eye:v2
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
      - ANALYTICS_MAX_BODY_BYTES=16384
      - ANALYTICS_MAX_PROPERTIES_BYTES=8192
      - ANALYTICS_MAX_SINGLE_PAGE_DURATION=4h
      - ANALYTICS_RATE_LIMIT_ENABLED=true
      - ANALYTICS_RATE_LIMIT_PER_MINUTE=120
      - ANALYTICS_RATE_LIMIT_BURST=240
      # Trust loopback and private network proxies by default.
      # Add CDN/public reverse-proxy ranges explicitly.
      - TRUSTED_PROXY_CIDRS=127.0.0.1/32,::1/128,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7
      - GRAPHQL_MAX_BODY_BYTES=1048576
      - GRAPHQL_MAX_COMPLEXITY=300
      - DASHBOARD_MAX_DAILY_RANGE_DAYS=730
      - DASHBOARD_MAX_HOURLY_RANGE_DAYS=31
      - DASHBOARD_MAX_FILTER_VALUES=100
      - DASHBOARD_MAX_FILTER_STRING_LENGTH=2048
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
    image: postgres:18.6-alpine
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
