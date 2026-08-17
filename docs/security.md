# Security and privacy model

This document describes the controls maintainers must preserve when changing authentication, the
dashboard transport, tracker collection, analytics identity, or dependencies.

## Trust boundaries and assets

| Boundary | Untrusted input | Protected asset | Primary controls |
| --- | --- | --- | --- |
| Dashboard GraphQL | JSON body, variables, cookies, forwarded headers | Admin session and site analytics | bounded body and operation complexity, auth mutation rate limit, HttpOnly scoped cookies, JWT validation, ownership checks, typed/masked errors |
| Tracker collection | site key, origin/referer, IP headers, user agent, JSON payload | analytics integrity and service availability | trusted-proxy resolution, configured-domain allowlist, two-stage rate limit, bounded body/properties, persistence-length validation, parameterized Bun queries |
| Dashboard rendering | stored paths, event data, and referrers | admin browser session | React escaping, HTTP(S)-only external referrer links, CSP, frame denial, no raw HTML sinks |
| GeoIP download | operator-controlled URL and license key | server filesystem and outbound network | explicit configuration ownership, bounded HTTP client behavior, archive validation, errors without client IP logging |
| SQLite/PostgreSQL | application-generated parameterized queries | users, sites, pseudonymous analytics | feature-owned persistence adapters, migrations, both-dialect tests, no transport exposure of row models |
| Dependency graph | registry packages and lifecycle scripts | build and release integrity | committed Bun/Go lock data, frozen install, Bun default-deny dependency scripts with one reviewed SWC exception, Bun audit and `govulncheck` |

Assets intentionally do not include a stable raw visitor identifier. Stored client hashes rotate
after a skipped UTC day and derive from a site-scoped keyed hash plus truncated IP prefix. Raw IP
addresses are used only in request processing and rate-limit memory; they are neither persisted nor
written to application logs.

## Required invariants

- Production sets distinct 32+ character `JWT_SECRET` and `ANALYTICS_IDENTITY_SECRET` values.
- Production uses HTTPS and leaves `SECURE_COOKIES=true`.
- Auth cookies remain HttpOnly, SameSite, runtime-`BASE_PATH` scoped, and uniquely named per mount.
- Every site-scoped GraphQL operation performs ownership authorization in the feature service.
- `POST /api/collect` accepts only a configured site domain and never reveals whether a site key or
  event definition exists.
- Collect payload limits mirror persistence limits: path/referrer 2,048 characters, UTM source and
  medium 128, UTM campaign 256. A request contains exactly one JSON value.
- Stored external URLs are data. The dashboard exposes a clickable referrer only after parsing an
  absolute `http:` or `https:` URL.
- Unexpected GraphQL errors are logged server-side and returned only as
  `INTERNAL_SERVER_ERROR: internal server error`.
- Raw passwords, JWTs, cookies, analytics identity secrets, and raw client IPs are never logged.
- The production server does not expose GraphiQL or another CDN-backed GraphQL playground.
- GraphQL operations above `GRAPHQL_MAX_COMPLEXITY` are rejected before resolver execution.
- Security headers remain compatible with the one-build runtime base-path matrix.

## Abuse cases covered

- Repeated login/register mutations are limited by trusted client IP.
- Rotating invalid site keys cannot bypass the collect IP admission limit.
- Cross-origin GraphQL mutations are rejected; collect CORS is derived from the site's configured
  domains rather than a wildcard.
- Oversized, trailing, or persistence-invalid collect payloads fail before feature/persistence work.
- Site IDs owned by another user return `FORBIDDEN`; missing resources remain distinct.
- `javascript:`, `data:`, scheme-relative, and malformed stored referrers are not rendered as links.
- Slow headers/bodies are bounded by `ReadHeaderTimeout`, `ReadTimeout`, body limits, and header size.
- GraphQL alias/selection amplification is bounded by gqlgen's fixed operation-complexity limit.

## Operational caveats

- If initial-admin credentials are absent, registration remains open until the first user claims the
  admin role. Public deployments must set `INITIAL_ADMIN_USERNAME` and `INITIAL_ADMIN_PASSWORD`
  before first start and leave `ALLOW_REGISTRATION` empty or false.
- Auth and collect rate limiters are process-local. A horizontally scaled deployment needs an
  approved shared limiter before assuming a global request budget.
- Refresh JWTs are stateless and cannot be individually revoked before expiry. Rotating
  `JWT_SECRET` invalidates every session and is the recovery action for suspected token theft.
- CSP currently allows inline scripts/styles because runtime base-path configuration and React style
  attributes are inline. Do not broaden other directives; a nonce-based runtime-config redesign
  requires its own approved plan and the full base-path browser matrix.
- `GEOIP_DOWNLOAD_URL` is operator-controlled. Do not derive it from dashboard or tracker input.

## Dependency audit status

Run from the existing Dev Container:

```sh
cd /workspaces/lovely-eye/dashboard
bun install --frozen-lockfile
bun audit
bun audit --prod

cd /workspaces/lovely-eye
task server:vuln
```

The 2026-08-17 review removed the temporary top-level security overrides for `@babel/core`,
`brace-expansion`, `immutable`, `js-yaml`, `nanoid`, `shell-quote`, `ws`, and `yaml`. Every parent
range and the frozen lockfile already resolve patched compatible versions. Keeping the overrides
made transitive tools appear to be project-owned dependencies and allowed Renovate to force
`js-yaml@5` across consumers that require `^4.1.0`. Production Bun dependencies and reachable Go
code have no reported vulnerabilities. No top-level override remains; a fresh lockfile also resolves
the direct and CLI GraphQL client preset requirements to the same `6.1.3` release.

The full development-tool audit retains one transitive advisory group with one high and one
moderate finding:

- `picomatch@2.3.1` through Vite, TanStack Router, shadcn, and GraphQL Codegen tool paths. Bun cannot
  apply its documented parent- or version-scoped override forms in Bun 1.3.14, while a global
  override would also force incompatible v4 consumers. A direct 2.3.2 pin does not replace the
  nested 2.3.1 edge. The affected glob input is repository-owned, not runtime/user input.

Verbose `govulncheck` also reports the module-only `GO-2026-5932` advisory for the unmaintained
`x/crypto/openpgp` package. Lovely Eye does not import that package; symbol and package results are
empty, and the module advisory has no fixed release.

Recheck this edge after any Bun, micromatch, or related tool release and no later than 2026-09-15.
Do not force a global override without generation, shadcn, build, and complete browser verification.

## Release verification

Before release, require all of the following:

1. `task check` and `task server:vuln` pass.
2. `bun audit --prod` reports no vulnerability; full `bun audit` matches only documented,
   unreachable dev-tool findings.
3. The 7-flow admin suite, `/`, `/lovely-eye`, `/tools/lovely-eye`, and same-origin multi-instance
   auth-isolation suites pass.
4. SQLite and PostgreSQL migration up/down/reapply tests pass in isolated databases.
5. `git diff --check` is clean and no `.env`, `.env.local`, `*.pem`, or `*.key` file is tracked.
