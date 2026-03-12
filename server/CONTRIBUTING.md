# Lovely Eye - Project Rules

## Authentication

- User email is optional, username is required
- If both `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` are set, the server creates that admin on startup
- If either initial-admin value is missing, the first self-registered user becomes admin
- `ALLOW_REGISTRATION` defaults to `false` when both initial-admin values are set, otherwise defaults to `true`
- `ALLOW_REGISTRATION=true` explicitly keeps registration open after the first user exists
- The first registration remains available whenever the database has no users

## API Structure

- **GraphQL API** (`/graphql`) - Contains all API methods for:
  - Authentication (register, login, refresh token)
  - Site management (create, update, delete, list)
  - Dashboard and analytics queries
  - User profile queries

- **REST API** - Limited to tracking functionality only:
  - `POST /api/collect` - Track page views and custom events
  - `GET /tracker.js` - Serve the tracking script

## Database

- Supports both SQLite and PostgreSQL
- SQLite is default for development (no configuration needed)
- `DB_DRIVER` and `DB_DSN` are optional - defaults to SQLite with `data/lovely_eye.db`
- To use PostgreSQL, set both `DB_DRIVER=postgres` and `DB_DSN=postgres://...`

## Analytics identity

- Visitor identity is server-generated and uses UTC-day-skipped rotation
- Identity is derived from a keyed hash of: site ID, truncated IP prefix (`/24` for IPv4, `/64` for IPv6), browser family, and device class
- The server checks today's and yesterday's hash; if only yesterday matches, it rewrites that client row to today's hash
- A new client is created only after a full UTC day was skipped
- Sessions still use 30-minute inactivity
- Country tracking stays separate from visitor identity and is only used for reporting when enabled
- Set `ANALYTICS_IDENTITY_SECRET` to control the identity key explicitly
- If `ANALYTICS_IDENTITY_SECRET` is unset, the server falls back to `JWT_SECRET`
- The dedicated identity secret helps reduce the impact of database-only leaks by making visitor IDs harder to reproduce

## Code structure
- [Migrations](./migrations/README.md)
- [E2E testing](./e2e/README.md)
- [Packages](./pkg/README.md)
- [App-related logic](./internal/README.md)

## Code Generation

- Run `task generate` after modifying `schema.graphqls` or e2e operations to regenerate GraphQL code
