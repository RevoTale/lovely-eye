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
- Feature behavior lives under `internal/<feature>`; database rows and Bun queries stay in that
  feature's `persistence` package.
- HTTP adapters live under `internal/transport/http`; configuration and database lifecycle live
  under `internal/platform`; `internal/app` is the explicit composition root.

## GraphQL errors

Resolver failures use GraphQL-native errors with a stable machine-readable `extensions.code`:

- `BAD_USER_INPUT` — invalid credentials, IDs, ranges, filters, or feature input.
- `UNAUTHENTICATED` — the operation requires a valid dashboard session.
- `FORBIDDEN` — the user is authenticated but policy or ownership denies the operation.
- `NOT_FOUND` — the requested resource does not exist.
- `CONFLICT` — a uniqueness or current-state conflict prevents the operation.
- `INTERNAL_SERVER_ERROR` — an unexpected failure; internal details are logged and never returned.

Add feature sentinel errors to `internal/graph/errors.go` when they cross the GraphQL boundary. UI
logic must branch on `extensions.code`, not on human-readable message text.

## Code Generation

- Keep feature SDL in `schema/<feature>.graphqls` and its handwritten transport adapter in
  `internal/graph/<feature>.resolvers.go`; generated gqlgen executor/model output stays isolated.
- Run `task generate` after modifying `schema/*.graphqls` or e2e operations to regenerate GraphQL code.
  A new resolver stub must be adopted as handwritten code before the freshness gate passes.
