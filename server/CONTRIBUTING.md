# Lovely Eye - Project Rules

## Authentication

- User email is optional, username is required
- First registered user automatically becomes admin
- Subsequent user registration is disabled by default
- Admin can enable registration via `ALLOW_REGISTRATION=true` environment variable
- Initial admin can be created via `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` environment variables

## API Structure

- **GraphQL API** (`/graphql`) - Contains all API methods for:
  - Authentication (register, login, refresh token)
  - Site management (create, update, delete, list)
  - Dashboard and analytics queries
  - User profile queries

- **REST API** - Limited to tracking functionality only:
  - `POST /api/collect` - Track page views
  - `POST /api/collect` - Track page views and custom events (legacy: `/api/event`)
  - `GET /tracker.js` - Serve the tracking script

## Database

- Supports both SQLite and PostgreSQL
- SQLite is default for development (no configuration needed)
- `DB_DRIVER` and `DB_DSN` are optional - defaults to SQLite with `data/lovely_eye.db`
- To use PostgreSQL, set both `DB_DRIVER=postgres` and `DB_DSN=postgres://...`'

## Code structure
- [Migrations](./migrations/README.md)
- [E2E testing](./e2e/README.md)
- [Packages](./pkg/README.md)
- [App-related logic](./internal/README.md)

## Code Generation

- Run `make generate` after modifying `schema.graphqls` or e2e operations to regenerate GraphQL code
