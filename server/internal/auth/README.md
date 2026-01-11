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

No CSRF tokens needed. See [discussion](https://www.reddit.com/r/node/comments/1im7yj0/comment/mc0ylfd/).

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | (random) | Secret key for signing tokens. Must be at least 32 characters. |
| `JWT_ACCESS_EXPIRY_MINUTES` | `15` | Access token lifetime in minutes |
| `JWT_REFRESH_DAYS` | `7` | Refresh token lifetime in days |
| `SECURE_COOKIES` | `false` | Set to `true` in production (requires HTTPS) |
| `COOKIE_DOMAIN` | (empty) | Cookie domain (leave empty for current domain) |
| `ALLOW_REGISTRATION` | `false` | Allow new user registration after first user |

## Roles

- `admin` - Full access (first user only)
- `user` - Site ownership required

Set initial admin via `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD` env vars.
