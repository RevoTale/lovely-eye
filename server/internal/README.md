# App internal logic
Here happend all the magic of the app. Future changes should follow the following rules:
- Every directory should try to be an independent, replacable module. 
- No Go-related files should be placed inside root of this directory.

A brief overview of the main modules:
- `./auth` - Authentication module with JWT-based auth using HTTP-only cookies. Handles user registration, login, token refresh, and CSRF protection.
- `./config` - Application configuration loader. Reads environment variables for server, database, and auth settings.
- `./database` - Database connection layer using [Bun ORM](https://github.com/uptrace/bun). Supports both SQLite and PostgreSQL.
- `./graph` - GraphQL API layer ([gqlgen](https://github.com/99designs/gqlgen)). Contains resolvers, generated code, and schema handlers for the dashboard.
- `./handlers` - HTTP handlers for REST endpoints. Currently handles analytics data collection (page views, events).
- `./middleware` - HTTP middleware (CORS, logging). Applied to HTTP routes for cross-cutting concerns.
- `./models` - Domain models with [Bun](https://github.com/uptrace/bun) annotations. Defines User, Site, Session, PageView, and Event entities.
- `./repository` - Data access layer. Provides CRUD operations for all models using [Bun ORM](https://github.com/uptrace/bun).
- `./server` - Application bootstrap and HTTP server setup. Wires all dependencies and configures routes.
- `./services` - Business logic layer. Contains SiteService and AnalyticsService with domain operations.