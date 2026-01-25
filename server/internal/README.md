# App internal logic

Rules:
- Each directory should be an independent, replaceable module
- No Go files in root of this directory

Modules:
- `./auth` - Authentication module with JWT-based auth using HTTP-only cookies. Handles user registration, login, token refresh, and CSRF mitigations via SameSite cookies.
- `./config` - Application configuration loader. Reads environment variables for server, database, and auth settings.
- `./database` - Database connection layer using [Bun ORM](https://github.com/uptrace/bun). Supports both SQLite and PostgreSQL.
- `./graph` - GraphQL API layer ([gqlgen](https://github.com/99designs/gqlgen)). Contains resolvers, generated code, and schema handlers for the dashboard.
- `./handlers` - HTTP handlers for REST endpoints. Currently handles analytics data collection (page views, events).
- `./middleware` - HTTP middleware (CORS, logging). Applied to HTTP routes for cross-cutting concerns.
- `./models` - Domain models with [Bun](https://github.com/uptrace/bun) annotations. Defines User, Site, Client, Session, Event, and event definition entities.
- `./repository` - Data access layer. Provides CRUD operations for all models using [Bun ORM](https://github.com/uptrace/bun).
- `./server` - Application bootstrap and HTTP server setup. Wires all dependencies and configures routes.
- `./services` - Business logic layer. Contains SiteService and AnalyticsService with domain operations.
