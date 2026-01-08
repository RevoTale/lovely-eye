# Authentication

This package provides authentication for Lovely Eye.

## Overview

Authentication is handled using JWT (JSON Web Tokens) stored in HTTP-only cookies. This approach combines the stateless benefits of JWTs with the security of HTTP-only cookies.

## How It Works

### Token Types

Two token types are used:

- **Access Token** - Short-lived (15 minutes by default). Used for authenticating API requests.
- **Refresh Token** - Long-lived (7 days by default). Used to obtain new access tokens without re-entering credentials.

### Authentication Flow

1. User submits username and password via the `login` mutation
2. Credentials are validated against the database (passwords are hashed with bcrypt)
3. On success, three cookies are set:
   - `le_access` - Access token (HTTP-only, not readable by JavaScript)
   - `le_refresh` - Refresh token (HTTP-only, not readable by JavaScript)
   - `le_csrf` - CSRF token (readable by JavaScript)
4. Subsequent requests are authenticated automatically via cookies
5. When the access token expires, it is refreshed automatically using the refresh token

### CSRF Protection

Cross-Site Request Forgery (CSRF) protection is implemented using the double-submit cookie pattern:

1. A CSRF token is stored in a cookie readable by JavaScript (`le_csrf`)
2. For state-changing requests (POST, PUT, DELETE), the token must be included in the `X-CSRF-Token` header
3. The server validates that the header value matches the cookie value

This prevents malicious sites from making authenticated requests on behalf of the user.

### Cookie Security

All authentication cookies use:

- `HttpOnly` flag (except CSRF cookie) - Prevents JavaScript access
- `SameSite=Strict` - Prevents cross-site request inclusion
- `Secure` flag (in production) - Only sent over HTTPS
- `Path=/` - Available for all routes

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | (random) | Secret key for signing tokens. Must be at least 32 characters. |
| `JWT_ACCESS_EXPIRY_MINUTES` | `15` | Access token lifetime in minutes |
| `JWT_REFRESH_DAYS` | `7` | Refresh token lifetime in days |
| `SECURE_COOKIES` | `false` | Set to `true` in production (requires HTTPS) |
| `COOKIE_DOMAIN` | (empty) | Cookie domain (leave empty for current domain) |
| `ALLOW_REGISTRATION` | `false` | Allow new user registration after first user |

## Usage

### Client-Side (JavaScript)

```javascript
// Login
const response = await fetch('/graphql', {
  method: 'POST',
  credentials: 'include', // Required for cookies
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    query: `mutation { login(input: { username: "admin", password: "secret" }) { user { id username } } }`
  })
});

// Subsequent requests - include CSRF token
function getCookie(name) {
  const value = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)');
  return value ? value.pop() : '';
}

const data = await fetch('/graphql', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json',
    'X-CSRF-Token': getCookie('le_csrf'),
  },
  body: JSON.stringify({
    query: `mutation { createSite(input: { domain: "example.com" }) { id } }`
  })
});

// Logout
await fetch('/graphql', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json',
    'X-CSRF-Token': getCookie('le_csrf'),
  },
  body: JSON.stringify({
    query: `mutation { logout }`
  })
});
```

### API Clients (Non-Browser)

For API clients that cannot use cookies, the `Authorization` header is supported as a fallback:

```bash
# Login and extract token from response
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { login(input: {username:\"admin\",password:\"secret\"}) { accessToken } }"}'

# Use token in subsequent requests
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"query":"{ me { id username } }"}'
```

## First User

The first user to register becomes an administrator. Subsequent registrations are disabled by default (controlled by `ALLOW_REGISTRATION`).

An initial admin can also be created via environment variables:

```bash
INITIAL_ADMIN_USERNAME=admin
INITIAL_ADMIN_PASSWORD=your-secure-password
```

## Roles

Two roles are supported:

- `admin` - Full access to all operations
- `user` - Limited access (site ownership required for most operations)

## Replacing the Implementation

To use a different authentication method (OAuth, sessions, etc.), implement the `auth.Service` interface and inject it during application startup.
