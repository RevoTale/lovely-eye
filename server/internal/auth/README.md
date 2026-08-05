# Authentication

JWT-based authentication with HttpOnly cookies.

## Tokens

- **Access Token** - 15 minutes. Authenticates API requests.
- **Refresh Token** - 7 days. Renews access tokens.

## Flow

1. User logs in with username/password
2. Credentials validated (bcrypt)
3. Two cookies set: `le_access` and `le_refresh` (HttpOnly)
4. Subsequent requests authenticated via cookies
5. Access token auto-refreshed when expired

## Cookie Settings

- `HttpOnly` - No JavaScript access (XSS protection)
- `SameSite=Strict` (production) or `Lax` (development) - CSRF protection
- `Secure` - HTTPS only in production
- `Path=/` - Available for all routes

GraphQL POST requests with a foreign `Origin` are rejected before resolvers run. SameSite cookies remain a second browser-side CSRF mitigation, and auth mutations are rate-limited per trusted client IP.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | generated at startup if empty | Secret key for signing tokens. Must be at least 32 characters when set. Set it explicitly in production if sessions should survive restarts. |
| `ANALYTICS_IDENTITY_SECRET` | falls back to `JWT_SECRET` | Optional dedicated secret for analytics visitor identity. Must be at least 32 characters when set. Analytics uses it for the daily UTC hashes behind UTC-day-skipped rotation, and it helps reduce the impact of database-only leaks by making visitor IDs harder to recompute. |
| `JWT_ACCESS_EXPIRY_MINUTES` | `15` | Access token lifetime in minutes |
| `JWT_REFRESH_DAYS` | `7` | Refresh token lifetime in days |
| `SECURE_COOKIES` | `true` | Set to `true` in production (requires HTTPS) |
| `COOKIE_DOMAIN` | (empty) | Cookie domain (leave empty for current domain) |
| `ALLOW_REGISTRATION` | `auto` | Post-bootstrap registration policy. Defaults to `false` when both `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` are set, otherwise defaults to `true`. The first registration is still available whenever no users exist. |
| `INITIAL_ADMIN_USERNAME` | (empty) | Optional initial admin username. Takes effect only when `INITIAL_ADMIN_PASSWORD` is also set. |
| `INITIAL_ADMIN_PASSWORD` | (empty) | Optional initial admin password. Takes effect only when `INITIAL_ADMIN_USERNAME` is also set. |
| `AUTH_RATE_LIMIT_ENABLED` | `true` | Enables per-process rate limiting for `login` and `register` GraphQL mutations. |
| `AUTH_RATE_LIMIT_ATTEMPTS` | `10` | Maximum auth mutation attempts per trusted client IP during the window. |
| `AUTH_RATE_LIMIT_WINDOW` | `15m` | Fixed auth rate-limit window. |

`ANALYTICS_IDENTITY_SECRET` is used by analytics tracking, not dashboard auth. It is documented here because it falls back to `JWT_SECRET` when unset.

## Roles

- `admin` - Full access (initial admin or first self-registered user when no initial admin is configured)
- `user` - Site ownership required

If both `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` are set, Lovely Eye creates that admin on startup and keeps registration disabled by default. If either value is missing, the first self-registered user becomes admin and registration stays enabled by default unless `ALLOW_REGISTRATION` is explicitly set.
