# App internal logic

Rules:
- Each directory should be an independent, replaceable module
- No Go files in root of this directory

Modules:
- `./auth` - Authentication module with JWT-based auth using HTTP-only cookies. Handles user registration, login, token refresh, credential validation, and cookie settings.
- `./config` - Application configuration loader. Reads environment variables for server, database, auth, auth rate limits, analytics identity, collect limits, single-page duration caps, trusted proxies, GraphQL body limits, and dashboard query caps.
- `./database` - Database connection layer using [Bun ORM](https://github.com/uptrace/bun). Supports both SQLite and PostgreSQL.
- `./graph` - GraphQL API layer ([gqlgen](https://github.com/99designs/gqlgen)). Contains resolvers, generated code, schema handlers, and dashboard date/filter limit enforcement.
- `./handlers` - HTTP handlers for REST endpoints. Handles analytics collection, including query `site_key` validation, body/property caps, collect rate limiting, domain validation, and trusted client-IP resolution.
- `./middleware` - HTTP middleware (CORS, auth rate limiting, logging). Applied to HTTP routes for cross-cutting concerns, including same-origin enforcement for dashboard GraphQL POSTs.
- `./models` - Domain models with [Bun](https://github.com/uptrace/bun) annotations. Defines User, Site, Client, Session, Event, and event definition entities.
- `./repository` - Data access layer. Provides CRUD operations and analytics aggregate queries using [Bun ORM](https://github.com/uptrace/bun). Pageview aggregates bucket by event time where accuracy requires it.
- `./server` - Application bootstrap and HTTP server setup. Wires all dependencies and configures routes.
- `./services` - Business logic layer. Contains SiteService and AnalyticsService with domain operations, including pseudonymous visitor identity with UTC-day-skipped rotation, server-side pageview/exit classification, 30-minute session handling, and capped single-page exit duration.
