# ADR-0001: Support an arbitrary runtime base path

## Status

Accepted

## Date

2026-08-05

## Context

Lovely Eye ships the dashboard as static files embedded in, or served by, the Go application. Operators may mount the same image at the origin root or below a reverse-proxy path that is unknown when the dashboard is built.

The deployment path can therefore be `/`, `/lovely-eye`, or a nested path such as `/tools/lovely-eye`. Building a separate dashboard for each mount path would couple the release artifact to deployment configuration. Absolute asset URLs such as `/assets/app.js` also escape a nested mount path, while prefixing individual application routes duplicates router responsibility and is easy to apply inconsistently.

## Decision

One dashboard build artifact must work unchanged at any normalized runtime `BASE_PATH`.

- Vite emits relative asset URLs by using a relative `base` (`./` or the equivalent empty value).
- The HTML base URL and `config.js` are resolved at runtime by the Go server.
- Runtime configuration is the source of truth for `BASE_PATH`, `API_URL`, and `GRAPHQL_URL`.
- The Go server mounts dashboard and API handlers below the configured base path and strips that prefix before static-file lookup.
- TanStack Router receives the normalized runtime base path through its `basepath` option. Route definitions, links, and navigation targets remain unprefixed.
- The static handler returns `index.html` for extensionless unknown paths so direct loads and refreshes of client routes work.
- Missing files with an extension return `404`; they must not fall back to HTML.
- Authentication cookie names are derived deterministically from the normalized runtime base path and their `Path` is scoped to that base path. This prevents login, refresh, or logout in one Lovely Eye instance from changing another instance mounted on the same origin.

No code may assume `/dashboard` or another fixed public mount path.

## Verification contract

Changes affecting the dashboard build, HTML, routing, runtime configuration, redirects, API URLs, or static server must verify one build artifact against all of these mounts:

| Runtime base path | Dashboard root | Nested client route | Prefixed asset and config | Unprefixed asset |
| --- | --- | --- | --- | --- |
| `/` | `200` | serves `index.html` | `200` | not applicable |
| `/lovely-eye` | `200` | serves `index.html` | `200` | `404` |
| `/tools/lovely-eye` | `200` | serves `index.html` | `200` | `404` |

At each non-root mount, generated links and navigation must retain the prefix, and API/GraphQL requests must use runtime-configured prefixed URLs. Test a trailing-slash entry and a direct refresh of at least one nested TanStack route.

The matrix must also run two instances on one origin under distinct base paths and verify independent login, refresh, and logout behavior. Cookies retain the accepted `HttpOnly`, `Secure`, and `SameSite` protections.

## Alternatives considered

### Build-time absolute base path

Rejected because it requires a separate build for each deployment path and prevents one container image from moving between mounts.

### Hardcoded `/dashboard`

Rejected because reverse proxies and operators control the public mount path.

### Manually prefix every route and URL

Rejected because prefix ownership becomes scattered across components and conflicts with TanStack Router's `basepath` handling.

## Consequences

- Dashboard assets must remain path-relative.
- Public URLs are deployment configuration, not build configuration.
- Router and server base-path normalization must stay aligned.
- Auth cookie naming, scoping, setting, reading, refreshing, and clearing must share the same normalized base-path policy.
- Base-path integration coverage is a required quality gate, not an optional deployment check.
